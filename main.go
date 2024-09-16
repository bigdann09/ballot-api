package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrz1836/go-sanitize"
	"github.com/robfig/cron"
)

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID          uint          `json:"id"`
	Source      Source        `json:"source"`
	Author      string        `json:"author"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	URL         string        `json:"url"`
	URLToImage  string        `json:"urlToImage"`
	PublishedAt string        `json:"publishedAt"`
	Content     template.HTML `json:"content"`
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

func main() {
	// initialize cron server
	c := StartCronScheduler()
	defer c.Stop()

	port := ":8001"
	fmt.Println("Server started at port ", port)

	// set mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/news", GetBallotNewsArticles)
	router.GET("/news/:slug", GetBallotNewsArticlesSlug)
	router.Run(port)
}

func GetBallotNewsArticles(c *gin.Context) {
	// get json data from file
	file, err := os.Open("news_articles.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var response Response
	json.Unmarshal(byteValue, &response)

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

func StartCronScheduler() *cron.Cron {
	// Create a new cron instance
	c := cron.New()

	// Add a cron job that runs every 10 seconds
	c.AddFunc("@every 04h00m00s", fetchNews)

	// Start the cron scheduler
	c.Start()
	return c
}

func fetchNews() {
	url := "https://newsapi.org/v2/everything?q=politics%20AND%20election&from=2024-09-14&apiKey=db10cebc16694fa99c5beb3c9eec64bf"

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
	for _, article := range data.Articles {
		if article.Content != "[Removed]" || article.Description != "[Removed]" {
			article.Slug = parseToSlug(article.Title)
			articles = append(articles, article)
		}
	}

	var response Response
	response.RefereshedAt = time.Now()
	response.Articles = articles

	parsedResponse, err := json.MarshalIndent(&response, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	// save to file
	if err := os.WriteFile("news_articles.json", parsedResponse, 0666); err != nil {
		log.Fatal(err)
	}
}

func parseToSlug(title string) string {
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, " ", "-")
	title = sanitize.Custom(title, "[^a-zA-Z0-9-]")
	return title
}
