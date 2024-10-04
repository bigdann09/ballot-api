package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mrz1836/go-sanitize"
	"github.com/robfig/cron/v3"
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
