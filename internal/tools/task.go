package tools

import (
	"io"

	"github.com/adepte-myao/test_parser/internal/models"
)

func WriteTasks(w io.Writer, tasks []models.Task) {
	for _, task := range tasks {
		w.Write([]byte(task.Question))

		w.Write([]byte("\n"))
		for _, option := range task.Options {
			w.Write([]byte("\t"))
			w.Write([]byte(option))
			w.Write([]byte("\n"))
		}

		w.Write([]byte("\tRight answer: "))
		w.Write([]byte(task.Answer))

		w.Write([]byte("\n\n"))
	}
}
