package controllers

import (
	"net/http"
	"github.com/ballot/internals/models"

	"github.com/gin-gonic/gin"
)

func GetNationalPolls(c *gin.Context) {
	polls := models.ScrapedPolls()
	c.JSON(http.StatusOK, polls)
}