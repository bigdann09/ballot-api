package models

import (
	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"
	"gorm.io/gorm"
)

type Referral struct {
	gorm.Model
	Referrer uint `json:"referrer"`
	Referee  uint `json:"referee"`
}

func NewReferral(referral *Referral) error {
	result := database.DB.Model(&Referral{}).Create(referral)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetReferrals(userID uint) ([]*utils.ReferralAPI, error) {
	var referrals []*utils.ReferralAPI
	result := database.DB.Model(&Referral{}).Where("referrer = ?", userID).Scan(&referrals)
	if result.Error != nil {
		return referrals, result.Error
	}

	if len(referrals) == 0 {
		referrals = []*utils.ReferralAPI{}
	}

	return referrals, nil
}
