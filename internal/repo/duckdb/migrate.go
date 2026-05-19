package duckdb

import (
	"database/sql"
	"fmt"
	"os"
)

func ApplyMigrationFile(db *sql.DB, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if _, err := db.Exec(string(b)); err != nil {
		return fmt.Errorf("apply migration %s: %w", path, err)
	}
	return nil
}
