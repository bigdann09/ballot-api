package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/mrz1836/go-sanitize"
	"github.com/robfig/cron/v3"
)

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID          uint   `json:"id"`
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

type Data struct {
	Status       string    `json:"status"`
	TotalResults uint      `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Response struct {
	RefereshedAt time.Time `json:"refereshedAt"`
	Articles     []Article `json:"articles"`
}

type Poll struct {
	Harris string `json:"harris"`
	Trump  string `json:"trump"`
}

type PollResult struct {
	Candidate string `json:"candidate"`
	Result    string `json:"result"`
}

func main() {
	// initialize cron server
	c := StartCronScheduler()
	defer c.Stop()

	port := ":8002"
	fmt.Println("Server started at port", port)

	// set mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(Cors())
	router.GET("/news", GetBallotNewsArticles)
	router.GET("/news/:slug", GetBallotNewsArticlesSlug)
	router.GET("/polls", GetNationalPolls)
	router.Run(port)
}

func GetBallotNewsArticles(c *gin.Context) {
	read, err := readFile("news_articles.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response Response
	json.Unmarshal(read, &response)

	// set content type and write response
	c.JSON(http.StatusOK, response)
}

func GetBallotNewsArticlesSlug(c *gin.Context) {
	slug := c.Param("slug")

	filebytes, err := os.ReadFile("news_articles.json")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response Response
	json.Unmarshal(filebytes, &response)

	// get article
	var article Article
	for _, a := range response.Articles {
		if strings.EqualFold(a.Slug, slug) {
			article = a
		}
	}

	// check if article is empty
	value := reflect.ValueOf(article)
	if value.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

func GetNationalPolls(c *gin.Context) {
	polls := scrapedPolls()
	c.JSON(http.StatusOK, polls)
}

func readFile(name string) ([]byte, error) {
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

// scrape polls
func scrapedPolls() []Poll {
	var polls = make([]Poll, 1)
	c := colly.NewCollector(
		colly.AllowedDomains("https://www.270towin.com", "www.270towin.com", "270towin.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.CacheDir("./polls_cache"),
	)

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	c.OnHTML("table#polls", func(e *colly.HTMLElement) {
		var poll Poll
		poll.Harris = e.ChildText("tr#poll_avg_row > td:nth-child(2)")
		poll.Trump = e.ChildText("tr#poll_avg_row > td:nth-child(3)")
		polls = append(polls, poll)
	})

	c.Visit("https://www.270towin.com/2024-presidential-election-polls/national")

	var modpolls []Poll
	for key, poll := range polls {
		if key == 1 {
			pollValue := reflect.ValueOf(poll)
			if !pollValue.IsZero() {
				modpolls = append(modpolls, poll)
			}
		}
	}

	return modpolls
}

// fetch news
func fetchNews() {
	url := "https://newsapi.org/v2/everything?q=politics%20AND%20election&from=%s&apiKey=db10cebc16694fa99c5beb3c9eec64bf"
	url = strings.Replace(url, "%s", formatDate(3), -1)

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

	// Add a cron job that runs every 10 seconds
	// c.AddFunc("@every 00h00m10s", scrapedPolls)

	// Start the cron scheduler
	c.Start()
	return c
}

// setup cors for api
func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"X-Requested-With", "Cache-Control", "Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}
