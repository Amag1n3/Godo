package task

import (
	"time"
)

type Task struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CreatedAt    string `json:"createdat"`
	DeadlineDate string `json:"deadlinedate"`
	DeadlineTime string `json:"deadlinetime"`
	Status       string `json:"status"`
}

func (task *Task) isOverdue() (bool, error) {
	deadlineStr := task.DeadlineDate + " " + task.DeadlineTime
	deadline, err := time.Parse("02/01/2006 1504", deadlineStr)
	if err != nil {
		return false, err
	}
	return time.Now().After(deadline), nil
}
