package models

import (
	_ "github.com/ballot/internals/database"
	_ "github.com/ballot/internals/utils"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name  string `json:"name"`
	Type  string `json:"type"`
	Link  string `json:"link" gorm:"default:''"`
	Point int64  `json:"point"`
	Validate bool `json:"validate" gorm:"default:false"`
}