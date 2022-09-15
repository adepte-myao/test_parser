package models

type Task struct {
	Question string
	Options  []string
	Answer   string
}

func NewTask() *Task {
	return &Task{}
}
