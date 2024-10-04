package app

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/models"
	"github.com/ballot/internals/utils"
)

func init() {
	// initialize cron server
	c := utils.StartCronScheduler()
	defer c.Stop()

	// load env variables
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	config := map[string]string{
		"host": os.Getenv("DB_HOST"),
		"port": os.Getenv("DB_PORT"),
		"user": os.Getenv("DB_USER"),
		"pass": os.Getenv("DB_PASS"),
		"name": os.Getenv("DB_NAME"),
		"ssl":  os.Getenv("SSL_MODE"),
	}

	// connect to database
	database.Connect(config)

	// auto create tables
	database.DB.AutoMigrate(&models.User{}, &models.Activity{}, &models.Point{}, &models.Referral{}, &models.Task{}, &models.UserTask{}, &models.Candidate{}, &models.Vote{})
}
