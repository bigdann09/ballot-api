package app

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/models"
)

func init() {

	isProd, _ := strconv.ParseBool(os.Getenv("PRODUCTION"))

	if !isProd {
		// load env variables
		err := godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}

	// set timezone
	os.Setenv("TZ", os.Getenv("TIMEZONE"))

	config := map[string]string{
		"host": os.Getenv("DB_HOST"),
		"port": os.Getenv("DB_PORT"),
		"user": os.Getenv("DB_USER"),
		"pass": os.Getenv("DB_PASS"),
		"name": os.Getenv("DB_NAME"),
		"ssl":  os.Getenv("SSL_MODE"),
	}

	// utils.PlaywrightArticleScraper()

	// connect to database
	database.Connect(config)

	// auto create tables
	database.DB.AutoMigrate(&models.User{}, &models.Activity{}, &models.Point{}, &models.Referral{}, &models.Task{}, &models.UserTask{}, &models.Candidate{}, &models.Vote{})
}
