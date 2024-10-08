package controllers

import (
	"net/http"
	"os"

	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"
	"github.com/gin-gonic/gin"
)

func MakeVoteController(c *gin.Context) {
	var vote utils.MakeVote
	if err := c.BindJSON(&vote); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
		return
	}

	user, err := utils.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// get user last vote activity
	activity, err := models.GetActivity(int64(user.TGID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// check if candidate with name exists
	if found := models.CheckCandidate(vote.Candidate); !found {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "candidate with name does not exist",
		})
		return
	}

	// fetch candidate
	candidate, err := models.GetCandidateByName(vote.Candidate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	if activity.LastVotedAt == "" {
		models.MakeVote(&models.Vote{
			UserID:      uint(user.ID),
			CandidateID: candidate.ID,
		})

		// update activity
		models.UpdateVoteActivity(uint(user.ID))

		// add vote points to user
		models.UpdateTaskPoint(uint(user.ID), uint64(utils.ParseStringToInt(os.Getenv("VOTE_POINTS"))))

	} else {
		last_voted := utils.CalculateTimeDiff(activity.LastVotedAt)
		if last_voted < 24 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  http.StatusTooManyRequests,
				"message": "You can only vote once every 24 hours",
			})
			return
		} else {
			models.MakeVote(&models.Vote{
				UserID:      uint(user.ID),
				CandidateID: candidate.ID,
			})

			// update activity
			models.UpdateVoteActivity(uint(user.ID))

			// add vote points to user
			models.UpdateTaskPoint(uint(user.ID), uint64(utils.ParseStringToInt(os.Getenv("VOTE_POINTS"))))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Vote casted successfully",
	})
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
	votes, err := models.GetDailyVotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, votes)
}
