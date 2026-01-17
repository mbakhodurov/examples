package migrator

import (
	"database/sql"

	"github.com/pressly/goose"
)

type Migrator struct {
	db          *sql.DB
	migratorDir string
}

func NewMigrator(db *sql.DB, migratorDir string) *Migrator {
	return &Migrator{
		db:          db,
		migratorDir: migratorDir,
	}
}

func (m *Migrator) Up() error {
	err := goose.Up(m.db, m.migratorDir)
	if err != nil {
		return err
	}
	return nil
}
