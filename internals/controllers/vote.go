package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"
	"github.com/gin-gonic/gin"
)

func MakeVoteController(c *gin.Context) {

}

func GetAllVotesController(c *gin.Context) {
	votes, err := models.GetAllVotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, votes)
}

func GetDailyVotesController(c *gin.Context) {
	// timee := time.Now().Sub(time.Date(2024, 10, 02, 23, 00, 45, 00, time.UTC))
	// old, _ := time.Parse("2006-01-02 03:04:05", "2024-10-03 11:15:39")
	// now := time.Now()
	// diff := now.Sub(old)

	daily, err := models.GetDailyVotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var data []*utils.VoteAPI
	for _, d := range daily {
		if d.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
			data = append(data, d)
		}
	}

	fmt.Println(data)

	c.JSON(http.StatusOK, data)
}
