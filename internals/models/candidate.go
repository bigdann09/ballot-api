package models

import (
	"time"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"

	"gorm.io/gorm"
)

type Candidate struct {
	gorm.Model
	Name       string    `json:"name"`
	Votes      uint      `json:"votes" gorm:"default:0"`
	LastVoteAt time.Time `json:"last_vote"`
}

func NewCandidate(name string) error {
	result := database.DB.Create(&Candidate{Name: name})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAllCandidates() ([]*utils.CandidateAPI, error) {
	var candidates []*utils.CandidateAPI
	result := database.DB.Model(&Candidate{}).Select("id", "name", "votes", "last_vote_at").Scan(&candidates)
	if result.Error != nil {
		return candidates, result.Error
	}

	if len(candidates) == 0 {
		candidates = []*utils.CandidateAPI{}
	}

	return candidates, nil
}

func GetCandidateByName(name string) (utils.CandidateAPI, error) {
	var candidate utils.CandidateAPI
	result := database.DB.Model(&Candidate{}).Where("name", name).Select("id", "name", "votes", "last_vote_at").Scan(&candidate)
	if result.Error != nil {
		return candidate, result.Error
	}
	return candidate, nil
}
