package models

import (
	_ "github.com/ballot/internals/database"
	_ "github.com/ballot/internals/utils"
	"gorm.io/gorm"
)

type Referral struct {
	gorm.Model
	Referrer uint `json:"referrer"`
	Referee uint `json:"referee"`
}