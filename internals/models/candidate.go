package models

import (
	"fmt"
	"time"

	"github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"

	"gorm.io/gorm"
)

type Candidate struct {
	gorm.Model
	Name       string    `json:"name"`
	Votes      uint      `json:"votes" gorm:"default:0"`
	Party      string    `json:"party"`
	LastVoteAt time.Time `json:"last_vote"`
}

func NewCandidate(name, party string) error {
	if !CheckCandidate(name) {
		result := database.DB.Create(&Candidate{
			Name:  name,
			Party: party,
		})
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
	return fmt.Errorf("candidate already exists")
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

func CheckCandidate(name string) bool {
	candidate, _ := GetCandidateByName(name)
	return candidate.ID != 0
}
