package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/ballot/internals/models"

	"github.com/gin-gonic/gin"
)

func GetBallotNewsArticles(c *gin.Context) {
	articles := models.GetArticles()
	if len(articles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No articles",
		})
		return
	}

	// set content type and write response
	c.JSON(http.StatusOK, articles)
}

func GetBallotNewsArticlesSlug(c *gin.Context) {
	slug := c.Param("slug")

	filebytes, err := os.ReadFile("articles.json")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var articles []*models.Article
	json.Unmarshal(filebytes, &articles)

	// get article
	var article *models.Article
	for _, a := range articles {
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
