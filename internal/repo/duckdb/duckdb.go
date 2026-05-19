package duckdb

import "database/sql"

type Adapter struct {
	DB *sql.DB
}

func New(db *sql.DB) *Adapter {
	return &Adapter{DB: db}
}
