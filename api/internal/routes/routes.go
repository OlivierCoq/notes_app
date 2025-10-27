package routes

import (
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/app"
	"github.com/go-chi/chi/v5"
)


func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	// Grouping routes and applying middleware can be done here if needed
	// the purpose of this r.Group method is to create a sub-router with specific middleware applied to it.
	// This is useful for applying middleware to a set of routes that share common requirements, such as authentication.
	// So all routes defined within this group will have the Authenticate middleware applied to them.
	// This helps keep the code organized and ensures that the middleware is consistently applied to all relevant routes.
	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate) // Apply the authentication middleware to all routes in this group

		// Note routes
		r.Get("/notes/{id}", app.Middleware.RequireUser(app.NoteHandler.HandleGetNoteByID))
		r.Get("/notes", app.Middleware.RequireUser(app.NoteHandler.HandleListNotesByUserID))
		r.Post("/notes", app.Middleware.RequireUser(app.NoteHandler.HandleCreateNote))
		r.Patch("/notes/{id}", app.Middleware.RequireUser(app.NoteHandler.HandleUpdateNote))
		r.Delete("/notes/{id}", app.Middleware.RequireUser(app.NoteHandler.HandleDeleteNote))
	})

	// Define routes and their handlers here
	r.Get("/health", app.HealthCheck) // Health check endpoint

	// // User registration route
	r.Post("/users/register", app.UserHandler.HandleRegisterUser)

	// // Token creation route
	// r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

	// // Logging user out:
	// r.Delete("/tokens/authentication", app.Middleware.RequireUser(app.TokenHandler.HandleRevokeToken))

	return r
}