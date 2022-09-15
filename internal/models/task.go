package models

type Task struct {
	Question string
	Options  []string
	Answer   string
	IsValid  bool
}

func NewTask() *Task {
	return &Task{
		IsValid: false,
	}
}
