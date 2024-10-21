package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept", "X-Requested-With", "Origin", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600 * 24 * 7 * 4 * 3,
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

	// NewsRoute
	news := router.Group("/api/news")
	news.GET("/user/:tg_id", controllers.GetBallotNewsArticles)
	news.GET("/:slug", controllers.GetBallotNewsArticlesSlug)

	// PollsRoute
	router.GET("/api/polls", controllers.GetNationalPolls)

	// CandidateRoute
	candidates := router.Group("api/candidates")
	candidates.Use(middleware.AuthMiddleware()) // apply middleware authentication
	candidates.GET("/", controllers.GetAllCandidatesController)
	candidates.GET("/:name", controllers.GetCandidateController)

	// UserRoute
	users := router.Group("api/users")
	users.Use(middleware.AuthMiddleware()) // apply middleware authentication
	users.GET("/leaderboard", controllers.GetLeaderboardsController)
	users.GET("/referrals", controllers.GetUserReferralsController)
	users.GET("/party", controllers.GetUserPartiesController)
	users.POST("/articles/claim", controllers.ClaimArticleRewardController)

	router.GET("/api/user", middleware.AuthMiddleware(), controllers.GetUserController)
	router.GET("/api/users/:tg_id", controllers.CheckUserController)
	router.POST("/api/users/onboard", controllers.OnboardUserController)

	// VoteRoute
	votes := router.Group("api/votes")
	votes.Use(middleware.AuthMiddleware()) // apply middleware authentication
	votes.GET("/", controllers.GetAllVotesController)
	votes.GET("/daily", controllers.GetDailyVotesController)
	votes.POST("/", controllers.MakeVoteController)
	votes.POST("/share/claim", controllers.ClaimVoteExtraRewardController)

	router.GET("/rrer", controllers.ActivityCheckerController)

	// TaskRoute
	tasks := router.Group("api/tasks")
	tasks.Use(middleware.AuthMiddleware()) // apply middleware authentication
	tasks.GET("/", controllers.GetAllTasksController)
	tasks.POST("/", controllers.StoreTaskController)
	tasks.PUT("/:uuid/complete", controllers.MarkTaskCompleteController)

	// JWTRoute
	jwt := router.Group("/api/jwt")
	jwt.GET("/:user_id/refresh", controllers.JWTRefreshController)
	jwt.POST("/verify", controllers.JWTVerifyController)
}

func Run(port string) {
	fmt.Println("started server at port", port)
	err := router.Run(port)
	fmt.Println(err)
}
