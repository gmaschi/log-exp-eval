package expstore

import "database/sql"

type Store interface {
	Querier
}

type exp struct {
	*Queries
	db *sql.DB
}

// NewStore returns a new Store interface
func NewStore(db *sql.DB) Store {
	store := &exp{
		Queries: New(db),
		db:      db,
	}

	return store
}
