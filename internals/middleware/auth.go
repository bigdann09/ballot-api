package middleware

import (
	"net/http"
	"time"

	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetTokenFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": err.Error(),
			})
			return
		}

		// validate token
		claims, err := utils.ValidateJWTToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  http.StatusUnprocessableEntity,
				"message": "Claims are not yet known",
			})
			return
		}

		expires_at := claims["expires_at"].(float64)
		if float64(time.Now().Unix()) > expires_at {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Token expired",
			})
			return
		}

		// get user
		issuer := claims["user"].(float64)
		user, _ := models.GetUser(int64(issuer))

		c.Set("user", user)
		c.Next()
	}
}
