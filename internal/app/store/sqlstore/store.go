package sqlstore

import (
	"database/sql"

	"github.com/ArtemGontar/betting/internal/app/store"
)

type Store struct {
	db                    *sql.DB
	matchResultRepository *MatchResultRepository
}

func New(databaseURL string) (*Store, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) MatchResult() store.MatchResultRepository {
	if s.matchResultRepository != nil {
		return s.matchResultRepository
	}
	return &MatchResultRepository{
		store: s,
	}
}
