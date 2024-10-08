package models

import (
	"fmt"

	"github.com/ballot/internals/database"
	_ "github.com/ballot/internals/database"
	"github.com/ballot/internals/utils"
	_ "github.com/ballot/internals/utils"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Link     string `json:"link" gorm:"default:''"`
	Point    int64  `json:"point"`
	Validate bool   `json:"validate" gorm:"default:false"`
	Duration string `json:"duration"`
}

func NewTask(task *utils.TaskCreateApiRequest) error {
	result := database.DB.Create(&Task{
		UUID:     utils.NewUUID(),
		Name:     task.Name,
		Type:     task.Type,
		Link:     task.Link,
		Point:    task.Point,
		Validate: task.Validate,
		Duration: task.Duration,
	})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("error creating task")
	}
	return nil
}

func GetAllTasks() ([]*utils.TaskAPI, error) {
	var tasks []*utils.TaskAPI
	result := database.DB.Model(&Task{}).Scan(&tasks)
	if result.Error != nil {
		return tasks, result.Error
	}
	return tasks, nil
}

func GetTask(taskID uint) (*utils.TaskAPI, error) {
	var task utils.TaskAPI
	result := database.DB.Model(&Task{}).Where("id = ?", taskID).Scan(&task)
	if result.Error != nil {
		return &task, result.Error
	}
	return &task, nil
}

func GetTaskByName(name string) (*utils.TaskAPI, error) {
	var task utils.TaskAPI
	result := database.DB.Model(&Task{}).Where("name = ?", name).Scan(&task)
	if result.Error != nil {
		return &task, result.Error
	}
	return &task, nil
}

func GetTaskByUUID(uuid string, tgID int64) (*utils.TaskAPI, error) {
	var task utils.TaskAPI
	result := database.DB.Model(&Task{}).Where("uuid = ?", uuid).Scan(&task)
	if result.Error != nil {
		return &task, result.Error
	}

	// check task
	task.Completed = CheckTask(tgID, task.ID)

	return &task, nil
}

func CheckTaskByName(name string) bool {
	task, _ := GetTaskByName(name)
	return task.ID != 0
}

func DeleteTask(uuid string) error {
	// result := task.Delete
	result := database.DB.Unscoped().Where("uuid = ?", uuid).Delete(&Task{})
	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return fmt.Errorf("error deleting task")
	}
	return nil
}
