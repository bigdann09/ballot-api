package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tgbotapi "gitlab.com/kingofsystem/telegram-bot-api/v5"

	"github.com/ballot/internals/bot"
	"github.com/ballot/internals/controllers"
	"github.com/ballot/internals/middleware"
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

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CandidateRoute
	candidates := router.Group("api/candidates")
	candidates.Use(middleware.AuthMiddleware()) // apply middleware authentication
	candidates.GET("/", controllers.GetAllCandidatesController)
	candidates.GET("/:name", controllers.GetCandidateController)

	// UserRoute
	users := router.Group("api/users")
	users.Use(middleware.AuthMiddleware()) // apply middleware authentication
	users.GET("/:tg_id", controllers.GetUserController)
	users.GET("/leaderboard", controllers.GetLeaderboardsController)
	users.GET("/referrals", controllers.GetUserReferralsController)
	users.GET("/party", controllers.GetUserPartiesController)
	router.POST("/api/users/onboard", controllers.OnboardUserController)

	// VoteRoute
	votes := router.Group("api/votes")
	votes.Use(middleware.AuthMiddleware()) // apply middleware authentication
	votes.GET("/", controllers.GetAllVotesController)
	votes.GET("/daily", controllers.GetDailyVotesController)
	votes.POST("/", controllers.MakeVoteController)

	// TaskRoute
	tasks := router.Group("api/tasks")
	tasks.Use(middleware.AuthMiddleware()) // apply middleware authentication
	tasks.GET("/", controllers.GetAllTasksController)
	tasks.POST("/", controllers.StoreTaskController)
	tasks.PUT("/:uuid/complete", controllers.MarkTaskCompleteController)

	referrals := router.Group("api/referrals")
	referrals.Use(middleware.AuthMiddleware()) // apply middleware authentication
	router.POST("api/referrals", controllers.StoreReferralController)
}

func Run(port string) {
	fmt.Println("started server at port", port)
	router.Run(port)
}
