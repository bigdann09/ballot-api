package controllers

import (
	"net/http"
	"strconv"

	"github.com/ballot/internals/models"

	"github.com/gin-gonic/gin"
)

func GetUserController(c *gin.Context) {
	tgID, _ := strconv.Atoi(c.Param("tg_id"))

	user, err := models.GetUser(int64(tgID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// get user point
	point, _ := models.GetPoint(user.ID)
	user.ReferralPoint = uint64(point.ReferralPoint)
	user.TaskPoint = uint64(point.TaskPoint)

	c.JSON(http.StatusOK, user)
}

func GetUserReferralsController(c *gin.Context) {
	tg_id, _ := strconv.Atoi(c.Param("tg_id"))

	// get user
	user, err := models.GetUser(int64(tg_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// get user referrals
	referrals, err := models.GetReferrals(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, referrals)
}

func GetLeaderboardsController(c *gin.Context) {
	// fetch all points
	user_points, err := models.GetAllPoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var points []map[string]interface{}
	for _, point := range user_points {
		total := point.ReferralPoint + point.TaskPoint
		user, _ := models.GetUserByID(uint(point.UserID))

		if total > 0 {
			points = append(points, map[string]interface{}{
				"tg_id":        user.TGID,
				"total_points": total,
			})
		}
	}

	if len(points) == 0 {
		points = []map[string]interface{}{}
	}

	c.JSON(http.StatusOK, points)
}
