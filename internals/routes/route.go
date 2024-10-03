package routes

import (
	"fmt"
	"encoding/json"
	"time"
	"net/http"
	"io"
	"os"

	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"

	"github.com/ballot/internals/bot"
	"github.com/ballot/internals/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

func Cors(origins ...string) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"X-Requested-With", "Cache-Control", "Origin", "Accept-Encoding", "X-CSRF-Token", "Content-Type", "Authorization", "Accept", "User-Agent", "Pragma"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
}

func RegisteredRoutes() {
	// BotRoute
	router.POST(fmt.Sprintf("/%s", os.Getenv("TELEGRAM_APIKEY")), func(c *gin.Context) {
		response, err := io.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		var update *tgbotapi.Update
		json.Unmarshal(response, &update)

		ballot_bot, err := bot.New(os.Getenv("TELEGRAM_APIKEY"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ballot_bot.HandleUpdates(update)
	})	

	// UserRoute
	user := router.Group("users")
	user.GET("/:tg_id", controllers.GetUserController)
	user.GET("/:tg_id/referrals", controllers.GetUserReferralsController)
	user.GET("/leaderboard", controllers.GetLeaderboardsController)
}

func Run(port string) {
	fmt.Println("started server at port", port)
	router.Run(port)
}