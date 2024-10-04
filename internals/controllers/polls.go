package controllers

import (
	"github.com/ballot/internals/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNationalPolls(c *gin.Context) {
	polls := models.ScrapedPolls()
	c.JSON(http.StatusOK, polls)
}
