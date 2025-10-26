package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/store"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/migrations"
)

// This is the main application struct that holds the dependencies for the app
type Application struct {
	Logger *log.Logger
	DB             *sql.DB // Add the database connection field
}

func NewApplication() (*Application, error) {

	// Create a new logger:
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime) 

	// Database connection
	pgDB, err := store.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}


	// Stores

	// Handlers

	// Midleware


	// Run database migrations using the embedded filesystem:
	// the "." means the current directory, which is where the migration files are located in the embedded FS
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		// panic and crash the app if migration fails:	
		panic(err)
	}



	app := &Application{
		Logger: logger,
		DB:     pgDB,
	}

	return app, nil
}