package controllers

import (
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"

	"github.com/gin-gonic/gin"
)

func GetUserReferralsController(c *gin.Context) {
	user, err := utils.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
		return
	}

	referrals, err := models.GetReferrals(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// fetch referee information
	referees := []*utils.UserAPI{}
	for _, referral := range referrals {
		referee, _ := models.GetUserByID(referral.Referee)
		point, _ := models.GetPoint(user.ID)
		referee.ReferralPoint = uint64(point.ReferralPoint)
		referee.TaskPoint = uint64(point.TaskPoint)

		referees = append(referees, referee)
	}

	c.JSON(http.StatusOK, referees)
}

func GetUserController(c *gin.Context) {
	user, err := utils.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
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

func GetLeaderboardsController(c *gin.Context) {

	authUser, err := utils.GetAuthUser(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
	}

	// fetch all points
	user_points, err := models.GetAllPoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var leaderboard []utils.LeaderboardAPI
	for _, point := range user_points {
		total := point.ReferralPoint + point.TaskPoint
		user, _ := models.GetUserByID(uint(point.UserID))

		if total > 0 {
			leaderboard = append(leaderboard, utils.LeaderboardAPI{
				TGID:        uint64(user.TGID),
				Username:    user.Username,
				TotalPoints: total,
			})
		}
	}

	if len(leaderboard) == 0 {
		leaderboard = []utils.LeaderboardAPI{}
	}

	sort.Sort(models.LeaderboardSort(leaderboard))

	var userPosition uint
	for key, lead := range leaderboard {
		if lead.TGID == uint64(authUser.TGID) {
			userPosition = uint(key + 1)
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"data":     leaderboard,
		"position": userPosition,
	})
}

func CheckUserController(c *gin.Context) {
	tgID := utils.ParseStringToInt(c.Param("tg_id"))

	var exists bool
	var status int

	// check if a user exist
	if found := models.CheckUser(int64(tgID)); found {
		status = http.StatusOK
		exists = found
	} else {
		status = http.StatusNotFound
		exists = found
	}

	c.JSON(status, gin.H{
		"status": status,
		"exists": exists,
	})

}

func OnboardUserController(c *gin.Context) {
	var user utils.NewUser
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
		return
	}

	if user.TGID == 0 || user.Username == "" || user.WalletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Missing required fields",
		})
		return
	}

	// check if user already exists
	if found := models.CheckUser(user.TGID); found {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "User already exists",
		})
		return
	}

	// add new user
	newUser, err := models.NewUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// check for referral
	ref := c.Query("ref")
	if ref != "" {
		referrer, err := models.GetUserByToken(ref)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if referrer.ID != 0 {
			// add new referral data
			err := models.NewReferral(&models.Referral{
				Referrer: referrer.ID,
				Referee:  newUser.ID,
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			// add points
			var points int
			if newUser.TGPremium {
				points = utils.ParseStringToInt(os.Getenv("REFERRAL_PREMIUM_POINTS"))
			} else {
				points = utils.ParseStringToInt(os.Getenv("REFERRAL_POINTS"))
			}

			// update referral points
			models.UpdateReferralPoint(referrer.ID, uint64(points))

			// send message to bot
			utils.SendMessageToBot(referrer.TGID, fmt.Sprintf("You just earned %d points for your referral", points))

			// update last login
			models.UpdateLoginActivity(newUser.ID)
		}
	}

	if newUser.TGPremium {
		models.UpdateTaskPoint(newUser.ID, uint64(utils.ParseStringToInt(os.Getenv("NEWUSER_PREMIUM_POINTS"))))
	} else {
		models.UpdateTaskPoint(newUser.ID, uint64(utils.ParseStringToInt(os.Getenv("NEWUSER_POINTS"))))
	}

	// set cookie
	token, err := utils.CreateJWTToken(user.TGID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	// Store cookie
	bearerToken := fmt.Sprintf("Bearer %s", token)
	c.SetSameSite(http.SameSiteNoneMode)
	c.Request.Header.Set("Authorization", bearerToken)

	// update last login
	models.UpdateLoginActivity(newUser.ID)

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"token":   token,
		"message": "User added successfully",
	})
}

func GetUserPartiesController(c *gin.Context) {
	republicans, _ := models.GetTotalUsersByParty("republican")
	democratics, _ := models.GetTotalUsersByParty("democratic")

	c.JSON(http.StatusOK, gin.H{
		"republicans": republicans,
		"democrats":   democratics,
	})
}
