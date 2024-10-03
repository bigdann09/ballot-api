package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(config map[string]string) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config["user"], config["pass"], config["host"], config["port"], config["name"])

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db
}
