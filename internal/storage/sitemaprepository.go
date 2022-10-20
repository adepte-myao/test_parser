package storage

import (
	"database/sql"

	"github.com/adepte-myao/test_parser/internal/models"
)

type SitemapRepository struct {
	store     *Store
	currentTx *sql.Tx
}

func NewSitemapRepository(store *Store) *SitemapRepository {
	return &SitemapRepository{
		store:     store,
		currentTx: nil,
	}
}

func (repo *SitemapRepository) GetAllTestLinks() ([]models.TestLink, error) {
	rows, err := repo.store.db.Query("SELECT test_id, link FROM tickets")
	if err != nil {
		return nil, err
	}

	testLinks := make([]models.TestLink, 0)
	for rows.Next() {
		var testLink models.TestLink
		err := rows.Scan(&testLink.TestId, &testLink.Link)
		if err != nil {
			return nil, err
		}

		testLinks = append(testLinks, testLink)
	}

	return testLinks, nil
}

func (repo *SitemapRepository) TruncateAllSitemapTables() error {
	_, err := repo.store.db.Exec("TRUNCATE sections CASCADE")
	return err
}

func (repo *SitemapRepository) CreateFilledSections(sections []models.Section) error {
	for _, section := range sections {
		var err error
		repo.currentTx, err = repo.store.db.Begin()
		if err != nil {
			repo.store.logger.Error("Cannot begin a transaction")
			return err
		}

		var sectionId int
		err = repo.currentTx.QueryRow(
			"INSERT INTO sections (name) VALUES ($1) RETURNING id",
			section.Name,
		).Scan(&sectionId)
		if err != nil {
			return err
		}

		err = repo.createFilledCertAreas(section.CertAreas, sectionId)
		if err != nil {
			repo.store.logger.Error(err.Error())
			repo.currentTx.Rollback()
			return err
		}

		if err = repo.currentTx.Commit(); err != nil {
			repo.store.logger.Error("Cannot commit a transaction")
			return err
		}
	}

	return nil
}

func (repo *SitemapRepository) createFilledCertAreas(certAreas []models.CertArea, sectionId int) error {
	for _, cerArea := range certAreas {
		var certAreaId int
		err := repo.currentTx.QueryRow(
			"INSERT INTO cert_area (section_id, name) VALUES ($1, $2) RETURNING id",
			sectionId,
			cerArea.Name,
		).Scan(&certAreaId)
		if err != nil {
			return err
		}

		err = repo.createFilledTests(cerArea.Tests, certAreaId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *SitemapRepository) createFilledTests(tests []models.Test, certAreaId int) error {
	for _, test := range tests {
		var testId int
		err := repo.currentTx.QueryRow(
			"INSERT INTO tests (cert_area_id, name) VALUES ($1, $2) RETURNING id",
			certAreaId,
			test.Name,
		).Scan(&testId)
		if err != nil {
			return err
		}

		repo.createTickets(test.TicketLinks, testId)
	}

	return nil
}

func (repo *SitemapRepository) createTickets(links []models.Link, testId int) error {
	for _, link := range links {
		// Not going to catch error there, because there will be lots of errors like
		// key value "lala" already exists
		// Hope there will be no other errors

		// Must check if there is the same links to avoid rollback
		var ticket_id int
		err := repo.currentTx.QueryRow(
			"SELECT id FROM tickets WHERE link = $1",
			string(link),
		).Scan(&ticket_id)
		if err == nil {
			continue
		}

		repo.currentTx.Exec(
			"INSERT INTO tickets (test_id, link) VALUES ($1, $2)",
			testId,
			string(link),
		)
	}

	return nil
}
