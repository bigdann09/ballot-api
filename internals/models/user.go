package models

import (
	"strings"

	"gorm.io/gorm"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"
)

type User struct {
	gorm.Model
	Username      string `json:"username" gorm:"default:''"`
	WalletAddress string `json:"wallet_address" gorm:"default:''"`
	TGID          int64  `json:"tg_id"`
	TGPremium     bool   `json:"tg_premium"`
	Token         string `json:"token"`
	Point         Point  `json:"point"`
	Party         string `json:"party" gorm:"default:''"`
}

func NewUser(user *utils.NewUser) (*User, error) {
	var newUser *User
	result := database.DB.Model(&User{}).Create(&User{
		Username:  user.Username,
		TGID:      user.TGID,
		TGPremium: user.TGPremium,
		Token:     utils.ReferralToken(10),
		Party:     strings.ToLower(user.Party),
	}).Scan(&newUser)

	if result.Error != nil {
		return newUser, result.Error
	}

	NewPoint(int(newUser.ID))
	NewActivity(int(newUser.ID))
	return newUser, nil
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

func GetTotalUsersByParty(party string) (int, error) {
	var total int
	result := database.DB.Model(&User{}).Select("COUNT(*) as total").Where("party = ?", party).Scan(&total)
	if result.Error != nil {
		return total, result.Error
	}
	return total, nil
}

func CheckUser(tgID int64) bool {
	user, _ := GetUser(tgID)
	return user.ID != 0
}
