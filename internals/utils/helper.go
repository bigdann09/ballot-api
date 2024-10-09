package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mrz1836/go-sanitize"
	"github.com/robfig/cron/v3"

	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

// generate from from news title
func parseToSlug(title string) string {
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, " ", "-")
	title = sanitize.Custom(title, "[^a-zA-Z0-9-]")
	return title
}

// format date for api
func formatDate(ago int) string {
	year, month, day := time.Now().Date()
	hour, min, sec := time.Now().Clock()
	prev_date := time.Date(year, month, (day - ago), hour, min, sec, 0, time.UTC)
	return prev_date.Format("2006-01-02")
}

// fetch news
func fetchNews() {
	url := "https://newsapi.org/v2/everything?q=politics%20AND%20election&from=%s&apiKey=db10cebc16694fa99c5beb3c9eec64bf"
	url = strings.Replace(url, "%s", formatDate(5), -1)

	// initialize new request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		os.Exit(1)
	}

	// send HTTP request
	res, err := client.Do(req)
	if err != nil {
		os.Exit(1)
	}

	// read response body
	parsed, err := io.ReadAll(res.Body)
	if err != nil {
		os.Exit(1)
	}
	defer res.Body.Close()

	var data Data
	json.Unmarshal(parsed, &data)

	// filter data
	var articles []Article
	count := 1
	for _, article := range data.Articles {
		if article.Content != "[Removed]" || article.Description != "[Removed]" {
			article.ID = uint(count)
			article.Slug = parseToSlug(article.Title)
			articles = append(articles, article)
			count++
		}
	}

	var response Response
	response.RefereshedAt = time.Now()
	response.Articles = articles

	parsedResponse, err := json.MarshalIndent(&response, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Last Fetched News at", time.Now().Format(time.RFC1123))

	err = os.WriteFile("news_articles.json", parsedResponse, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func ScrapeNewsArticle() {
	c := colly.NewCollector(
		colly.AllowedDomains("https://edition.cnn.com", "www.edition.cnn.com", "edition.cnn.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.CacheDir("./news_cache"),
	)

	c.SetRequestTimeout(120 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	c.OnHTML("article.live-story-post.liveStoryPost", func(e *colly.HTMLElement) {
		title := e.ChildText("h2.live-story-post__headline.inline-placeholder > strong")
		publishedBy := e.ChildText("span.live-story-post__byline.inline-placeholder")
		thumb := e.ChildAttr("div.image_live-story.image_live-story__hide-placeholder", "data-url")
		metadata := e.ChildText("span[data-editable]")
		e.ForEachWithBreak("p.paragraph.inline-placeholder.vossi-paragraph", func(i int, f *colly.HTMLElement) bool {
			return true
		})
		fmt.Println(title, publishedBy, thumb, metadata, metadata, contents)
	})

	c.Visit("https://edition.cnn.com/politics/live-news/trump-harris-election-10-08-24/index.html")
}

// configure cron
func StartCronScheduler() *cron.Cron {
	// Create a new cron instance
	c := cron.New(
		cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)),
		cron.WithLocation(time.UTC),
	)

	// Add a cron job that runs every 4 hours
	c.AddFunc("@every 04h00m00s", fetchNews)

	// Start the cron scheduler
	c.Start()
	return c
}

func ReadFile(name string) ([]byte, error) {
	// check if file exists
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return []byte{}, err
	}

	bbyte, err := os.ReadFile(name)
	if err != nil {
		return bbyte, err
	}

	return bbyte, nil
}

func ReferralToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func ValidateData(s interface{}, omittedFields []string) []map[string]string {
	structType := reflect.TypeOf(s)
	structValue := reflect.ValueOf(s)
	fieldNum := structValue.NumField()

	var data []map[string]string
	for i := 0; i < fieldNum; i++ {
		field := structValue.Field(i)
		fieldName := structType.Field(i).Name

		isSet := field.IsValid() && !field.IsZero()
		if !isSet && !slices.Contains(omittedFields, fieldName) {
			data = append(data, map[string]string{
				strings.ToLower(fieldName): "field value is required",
			})
		}
	}

	return data
}

func ParseStringToInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func CreateJWTToken(issuer int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"user":       issuer,
		"expires_at": time.Now().AddDate(2024, 12, 12).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func parseJWTToken(token string) string {
	if strings.Contains(token, "Bearer") {
		return strings.Split(token, " ")[1]
	}
	return token
}

func ValidateJWTToken(tokenString string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_KEY")), nil
	})

	if err != nil {
		return claims, err
	}

	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return claims, fmt.Errorf("Error retrieving claims")
	}

	return claims, nil
}

func GetAuthUser(c *gin.Context) (*UserAPI, error) {
	var empty *UserAPI
	data, ok := c.Get("user")
	if !ok {
		return empty, fmt.Errorf("unauthorized access")
	}

	user := data.(*UserAPI)
	return user, nil
}

func SendMessageToBot(tgID int64, message string) error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APIKEY"))
	if err != nil {
		return err
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp("Vote", tgbotapi.WebAppInfo{
				URL: MINIAPP_URL,
			}),
		),
	)

	msg := tgbotapi.NewMessage(tgID, message)
	msg.ReplyMarkup = markup
	_, err = bot.Send(msg)
	return err
}

func CurrentTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func CalculateTimeDiff(timeString string) int {
	parsed, _ := time.Parse(time.RFC3339, timeString)
	diff := time.Now().Sub(parsed).Hours()
	return int(math.Floor(diff))
}

func NewUUID() string {
	id := uuid.New()
	return id.String()
}

func GetTokenFromHeader(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("Authorization")
	token = parseJWTToken(token)
	if token == "" {
		return "", fmt.Errorf("Bearer token not found")
	}
	return token, nil
}
