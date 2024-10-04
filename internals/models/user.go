package models

import (
	"gorm.io/gorm"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"
)

type User struct {
	gorm.Model
	Username      string `json:"username" gorm:"default:''"`
	WalletAddress string `json:"username" gorm:"default:''"`
	TGID          int64  `json:"tg_id"`
	TGPremium     bool   `json:"tg_premium"`
	Token         string `json:"token"`
	Point         Point  `json:"point"`
}

func NewUser(user *User) (*User, error) {
	result := database.DB.Model(&User{}).Create(user)
	if result.Error != nil {
		return &User{}, result.Error
	}

	NewPoint(&Point{UserID: int(user.ID)})
	return user, nil
}

func GetUser(tgID int64) (*utils.UserAPI, error) {
	var user utils.UserAPI
	result := database.DB.Model(&User{}).Where("tg_id = ?", tgID).Select("id, username, tg_premium, tg_id, wallet_address, token").Scan(&user)
	if result != nil {
		return &user, result.Error
	}
	return &user, nil
}

func GetUserByID(ID uint) (*utils.UserAPI, error) {
	var user utils.UserAPI
	result := database.DB.Model(&User{}).Where("id = ?", ID).Select("id, username, tg_premium, tg_id, wallet_address").Scan(&user)
	if result != nil {
		return &user, result.Error
	}
	return &user, nil
}

func GetUserByToken(token string) (*utils.UserAPI, error) {
	var user utils.UserAPI
	result := database.DB.Model(&User{}).Where("token = ?", token).Select("id, username, tg_premium, tg_id, wallet_address").Scan(&user)
	if result != nil {
		return &user, result.Error
	}
	return &user, nil
}

func CheckUser(tgID int64) bool {
	user, _ := GetUser(tgID)
	return user.ID != 0
}
