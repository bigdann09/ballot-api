package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mrz1836/go-sanitize"

	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"
)

// generate from from news title
func ParseToSlug(title string) string {
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

func NextVoteTime() time.Time {
	year, month, day := time.Now().Date()
	nextVote := time.Date(year, month, day+1, 13, 0, 0, 0, time.Local)
	return nextVote
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
