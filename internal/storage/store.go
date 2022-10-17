package storage

import (
	"database/sql"
	"github.com/adepte-myao/test_parser/internal/config"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Store struct {
	config *config.StoreConfig
	db     *sql.DB
	logger *logrus.Logger
}

func NewStore(config *config.StoreConfig, logger *logrus.Logger) *Store {
	return &Store{
		config: config,
		logger: logger,
	}
}

func (s *Store) Open() error {
	s.logger.Info("Connecting to database: first attempt")

	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		// Sometimes database isn't up when we try to connect it.
		// Take some time, maybe db will be up
		s.logger.Info("Connecting to database: second attempt")

		time.Sleep(10 * time.Second)
		if db, err = sql.Open("postgres", s.config.DatabaseURL); err != nil {
			return err
		}
	}
	s.logger.Info("Connected to database")

	if err := db.Ping(); err != nil {
		return err
	}
	s.logger.Info("Ping to database is successful")

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}
