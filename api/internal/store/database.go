package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)


func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	fmt.Println("Database connection established! 🚀")
	// db.SetMaxOpenConns(), db.SetMaxIdleConns(), and db.SetConnMaxIdleTime()
	return db, nil
}

// MigrateFS applies database migrations from the provided fs.FS (embedded filesystem).
// It sets the base filesystem for goose to the provided migrationFS, runs the migrations,
// and then resets the base filesystem to nil.
func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)

	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {

	// Set the dialect for goose to PostgreSQL
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Apply all up migrations
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}