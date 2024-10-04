package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"
	"github.com/gin-gonic/gin"
)

func StoreReferralController(c *gin.Context) {
	var referral utils.ReferralCreateApiRequest
	if err := c.BindJSON(&referral); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check token if it exists
	referrer, err := models.GetUserByToken(referral.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check if user already exists
	if models.CheckUser(referral.TGID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User already exists",
		})
		return
	}

	// add new user
	user, _ := models.NewUser(&utils.NewUser{
		TGID:      referral.TGID,
		Username:  referral.Username,
		TGPremium: referral.TGPremium,
	})

	if referrer.ID != 0 {
		// add new referral data
		err := models.NewReferral(&models.Referral{
			Referrer: referrer.ID,
			Referee:  user.ID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// add points
		var points int
		if referral.TGPremium {
			points = utils.ParseStringToInt(os.Getenv("REFERRAL_PREMIUM_POINTS"))
		} else {
			points = utils.ParseStringToInt(os.Getenv("REFERRAL_POINTS"))
		}

		// update referral points
		models.UpdateReferralPoint(referrer.ID, uint64(points))

		// send message to bot
		utils.SendMessageToBot(referrer.TGID, fmt.Sprintf("You just earned %d points for your referral", points))

		// update last login
		models.UpdateLoginActivity(user.ID)

		c.JSON(http.StatusOK, gin.H{
			"message": "Referral data added successfully",
		})
	}
}
