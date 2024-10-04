package models

import (
	"time"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"
	"gorm.io/gorm"
)

type Vote struct {
	gorm.Model
	UserID      uint `json:"user_id"`
	CandidateID uint `json:"candidate_id"`
}

func MakeVote(vote *Vote) error {
	return database.DB.Model(&Vote{}).Create(vote).Error
}

func GetAllVotes() ([]*utils.VoteAPI, error) {
	var votes []*utils.VoteAPI
	result := database.DB.Model(&Vote{}).Select("id", "user_id", "candidate_id").Scan(&votes)
	if result.Error != nil {
		return votes, result.Error
	}

	if len(votes) == 0 {
		votes = []*utils.VoteAPI{}
	}

	return votes, nil
}

func GetDailyVotes() ([]*utils.VoteAPI, error) {
	var data []*utils.VoteAPI
	result := database.DB.Raw("SELECT id, user_id, candidate_id, DATE_TRUNC('day', created_at) as created_at FROM votes").Scan(&data)
	if result.Error != nil {
		return data, result.Error
	}

	var votes []*utils.VoteAPI
	if len(data) > 0 {
		for _, vote := range data {
			if vote.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
				votes = append(votes, vote)
			}
		}
	} else if len(data) == 0 {
		data = []*utils.VoteAPI{}
	}

	if len(votes) == 0 {
		votes = []*utils.VoteAPI{}
	}

	return votes, nil
}
