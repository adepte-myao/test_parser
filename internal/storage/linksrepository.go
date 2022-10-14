package storage

import (
	"github.com/adepte-myao/test_parser/internal/models"
)

type LinkRepository struct {
	store *Store
}

func NewLinksRepository(store *Store) *LinkRepository {
	return &LinkRepository{
		store: store,
	}
}

func (repo *LinkRepository) CreateRange(links []models.Link) error {
	tx, err := repo.store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, link := range links {
		_, err := tx.Exec(
			"INSERT INTO links (reference) VALUES ($1) RETURNING id",
			link,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *LinkRepository) GetAllLinks() ([]models.Link, error) {
	rows, err := repo.store.db.Query("SELECT reference FROM links")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]models.Link, 0)
	for rows.Next() {
		var link models.Link
		err := rows.Scan(&link)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	return links, nil
}

func (repo *LinkRepository) DeleteAll() error {
	_, err := repo.store.db.Exec("TRUNCATE links")
	if err != nil {
		return err
	}
	return nil
}
