package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/api"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/middleware"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/store"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/migrations"
)

// This is the main application struct that holds the dependencies for the app
type Application struct {
	Logger 			 *log.Logger
	DB           *sql.DB // Add the database connection field
	UserHandler  *api.UserHandler
	TokenHandler *api.TokenHandler
	Middleware   *middleware.UserMiddleware
	NoteHandler  *api.NoteHandler
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
	notesStore := store.NewPostgresNoteStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// Handlers
	noteHandler := api.NewNoteHandler(notesStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	// Middleware
	middlewareHandler := &middleware.UserMiddleware{
		UserStore: userStore,
	}
	userMiddleware := middlewareHandler

	// Run database migrations using the embedded filesystem:
	// the "." means the current directory, which is where the migration files are located in the embedded FS
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		// panic and crash the app if migration fails:	
		panic(err)
	}



	app := &Application{
		Logger:       logger,
		DB:           pgDB,
		UserHandler:  userHandler,
		TokenHandler: tokenHandler,
		Middleware:   userMiddleware,
		NoteHandler:  noteHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	/*
		- Purpose: To verify that the server is running and responsive.
		- Called by client (UI, some frontent)
		- needs 2 arguments: ResponseWriter and Request
		- ResponseWriter: used to send a response back to the client
		- Request: contains all the information about the incoming HTTP request. This is a pointer because it can be large and we want to avoid copying it. We
		also need it to persist and modify it, especially when dealing with middleware or request body.
		- In a real-world scenario, you might want to include more detailed health information,
		  such as database connectivity, external service status, etc.
		- In this example, we simply write "OK" to the response with a 200 status code.
	*/
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Status is available. A okay! 🟢\n")
}