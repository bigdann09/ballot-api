package models

import (
	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"
	"gorm.io/gorm"
)

type Activity struct {
	gorm.Model
	UserID       int    `json:"user_id"`
	LastLogin    string `json:"last_login" gorm:"default:null"`
	LastActivity string `json:"last_activity" gorm:"default:null"`
	LastVotedAt  string `json:"last_voted_at" gorm:"default:null"`
}

func NewActivity(userID int) error {
	result := database.DB.Model(&Activity{}).Create(&Activity{UserID: userID})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetActivity(userID int64) (*Activity, error) {
	var activity *Activity
	result := database.DB.Model(&Activity{}).Scan(&activity)
	if result.Error != nil {
		return activity, result.Error
	}

	return activity, nil
}

func UpdateVoteActivity(userID uint) error {
	result := database.DB.Model(&Activity{}).Where("user_id =?", userID).Updates(&Activity{
		LastActivity: utils.CurrentTime(),
		LastVotedAt:  utils.CurrentTime(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateLoginActivity(userID uint) error {
	result := database.DB.Model(&Activity{}).Where("user_id =?", userID).Updates(&Activity{
		LastActivity: utils.CurrentTime(),
		LastLogin:    utils.CurrentTime(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateLastActivity(userID uint) error {
	result := database.DB.Model(&Activity{}).Where("user_id =?", userID).Updates(&Activity{
		LastActivity: utils.CurrentTime(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
