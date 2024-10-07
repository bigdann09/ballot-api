package controllers

import (
	"net/http"
	"time"

	"github.com/ballot/internals/utils"
	"github.com/gin-gonic/gin"
)

func JWTRefreshController(c *gin.Context) {
	issuer := utils.ParseStringToInt(c.Param("user_id"))
	cookie, err := c.Cookie("Authorization")

	// parse token
	token := utils.ParseJWTToken(cookie)

	// validate token
	claims, err := utils.ValidateJWTToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": "Claims are not yet known",
		})
		return
	}

	if issuer != int(claims["user"].(float64)) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": http.StatusText(http.StatusUnauthorized),
		})
		return
	}

	expires_at := claims["expires_at"].(float64)
	if float64(time.Now().Unix()) > expires_at {
		new_token, err := utils.CreateJWTToken(int64(issuer))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"token":  new_token,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  token,
	})
}

func JWTVerifyController(c *gin.Context) {
	var jwt map[string]string
	if err := c.BindJSON(&jwt); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err.Error(),
		})
		return
	}

	token := jwt["token"]
	_, err := utils.ValidateJWTToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": true,
		})
	}
}
