package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	Source       Source        `json:"source"`
	Author       string        `json:"author"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	URL          string        `json:"url"`
	URLToImage   string        `json:"urlToImage"`
	PublishedAt  string        `json:"publishedAt"`
	Content      template.HTML `json:"content"`
	RefereshedAt time.Time     `json:"refereshedAt"`
}

type Data struct {
	Status       string    `json:"status"`
	TotalResults uint      `json:"totalResults"`
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

	var articles []Article
	json.Unmarshal(byteValue, &articles)

	// set content type and write response
	c.JSON(http.StatusOK, articles)
}

func StartCronScheduler() *cron.Cron {
	// Create a new cron instance
	c := cron.New()

	// Add a cron job that runs every 10 seconds
	c.AddFunc("@every 12h00m00s", fetchNews)

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
			article.RefereshedAt = time.Now()
			articles = append(articles, article)
		}
	}

	parsedArticles, err := json.MarshalIndent(&articles, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	// save to file
	if err := os.WriteFile("news_articles.json", parsedArticles, 0666); err != nil {
		log.Fatal(err)
	}
}
