package models

type UserTask struct {
	TGID   int64 `json:"tg_id"`
	TaskID uint  `json:"task_id"`
}