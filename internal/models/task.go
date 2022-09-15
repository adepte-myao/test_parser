package models

type Task struct {
	question string
	options  []string
	answer   string
	isValid  bool
}

func NewTask() *Task {
	return &Task{
		isValid: false,
	}
}
