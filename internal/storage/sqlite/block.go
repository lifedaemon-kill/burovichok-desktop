package sqlite

import "github.com/jmoiron/sqlx"

// TODO достать все, добавить строку
func NewBlockRepository(db *sqlx.DB) (BlocksStorage, error) {
	return &sqlite{
		DB: db,
	}, nil
}
