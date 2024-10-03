package controllers

import (
	"os"
	"reflect"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ballot/internals/utils"

	"github.com/gin-gonic/gin"
)

func GetBallotNewsArticles(c *gin.Context) {
	read, err := utils.ReadFile("news_articles.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response utils.Response
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

	var response utils.Response
	json.Unmarshal(filebytes, &response)

	// get article
	var article utils.Article
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