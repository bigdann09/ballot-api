package controllers

import (
	"net/http"

	"github.com/ballot/internals/models"
	"github.com/gin-gonic/gin"
)

func GetAllCandidatesController(c *gin.Context) {
	candidates, err := models.GetAllCandidates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, candidates)
}

func GetCandidateController(c *gin.Context) {
	name := c.Param("name")

	candidate, err := models.GetCandidateByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if candidate.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Candidate not found",
		})
		return
	}

	c.JSON(http.StatusOK, candidate)
}
