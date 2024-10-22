package models

import (
	"github.com/ballot/internals/database"
)

type UserTask struct {
	TGID   int64 `json:"tg_id"`
	TaskID uint  `json:"task_id"`
}

func CompleteTask(tgID int64, taskID uint) error {
	result := database.DB.Model(&UserTask{}).Create(&UserTask{
		TGID:   tgID,
		TaskID: taskID,
	})
	if result.Error != nil {
		return result.Error
	}

	// store user point

	return nil
}

func GetCompletedTask(tgID int64, taskID uint) (*UserTask, error) {
	var completed UserTask
	result := database.DB.Model(&UserTask{}).Where("tg_id = ? AND task_id = ?", tgID, taskID).Scan(&completed)
	if result.Error != nil {
		return &completed, result.Error
	}
	return &completed, nil
}

func CheckTask(tgID int64, taskID uint) bool {
	task, _ := GetCompletedTask(tgID, taskID)
	return task.TGID != 0 && task.TaskID != 0
}
