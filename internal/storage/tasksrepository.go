package storage

import "github.com/adepte-myao/test_parser/internal/models"

type TaskRepository struct {
	store *Store
}

func NewTaskRepository(store *Store) *TaskRepository {
	return &TaskRepository{
		store: store,
	}
}

func (repo *TaskRepository) CreateTask(task models.Task) error {
	var taskId int32
	repo.store.db.QueryRow(
		"INSERT INTO tasks (question, answer) VALUES ($1, $2) RETURNING id",
		task.Question,
		task.Answer,
	).Scan(&taskId)

	for _, option := range task.Options {
		repo.store.db.Exec(
			"INSERT INTO options (question_id, answer_option) VALUES ($1, $2)",
			taskId,
			option,
		)
	}

	return nil
}

func (repo *TaskRepository) DeleteAll() error {
	_, err := repo.store.db.Exec("TRUNCATE tasks")
	if err != nil {
		return err
	}
	return nil
}
