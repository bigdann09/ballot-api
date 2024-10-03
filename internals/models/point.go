package models

import (
	"fmt"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"
	"gorm.io/gorm"
)

type Point struct {
	gorm.Model
	UserID        int
	ReferralPoint uint64 `gorm:"default:0"`
	TaskPoint     uint64 `gorm:"default:0"`
}

func NewPoint(point *Point) error {
	result := database.DB.Create(point)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetPoint(userID uint) (*utils.PointAPI, error) {
	var point utils.PointAPI
	result := database.DB.Model(&Point{}).Where("user_id = ?", userID).Scan(&point)
	if result.Error != nil {
		return &point, result.Error
	}

	return &point, nil
}

func GetAllPoints() ([]*utils.PointAPI, error) {
	var points []*utils.PointAPI
	result := database.DB.Model(&Point{}).Order("referral_point DESC").Order("task_point").Scan(&points)
	if result.Error != nil {
		return points, result.Error
	}

	return points, nil
}

func UpdateReferralPoint(userID uint, referral_point uint64) error {
	point, _ := GetPoint(userID)
	point.ReferralPoint += referral_point

	result := database.DB.Model(&Point{}).Where("user_id = ?", userID).Updates(&Point{
		ReferralPoint: point.ReferralPoint,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("error updating user's referral point")
	}
	return nil
}
