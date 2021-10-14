package sqlstore

import (
	"database/sql"

	"github.com/ArtemGontar/betting/internal/app/store"
)

type Store struct {
	db                    *sql.DB
	matchResultRepository *MatchResultRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) MatchResult() store.MatchResultRepository {
	if s.matchResultRepository != nil {
		return s.matchResultRepository
	}
	return &MatchResultRepository{
		store: s,
	}
}
