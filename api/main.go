package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/app"
	"github.com/OlivierCoq/notes_app/api/notes_app_api/internal/routes"
)

// To learn more, visit: https://github.com/OlivierCoq/go_api_template


func main() { 

	var port int

	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()


	// Initialize the application (taken from internal/app/app.go):
	app, err := app.NewApplication()
	if err != nil {
		// Worst case scenario, we panic here with the error. Will crash the app
		panic(err)
	}
	// Log that the app has started
	app.Logger.Println("Notes app api started. Werk it! 🚀")

	// Database connection:
defer app.DB.Close()


// Routes and Handlers setup

		// Using chi for routing
	r := routes.SetupRoutes(app)

	// Use the chi router as the main handler for incoming requests
	http.Handle("/", r)


	// declare a new server with specific configurations
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port), // returns variable port as a string with a colon in front of it
		IdleTimeout:  time.Minute,              // how long to wait before closing idle connections
		ReadTimeout:  10 * time.Second,         // max duration for reading the entire request, including the body
		WriteTimeout: 30 * time.Second,         // max duration before timing out writes of the response
	}


// Start the server
	app.Logger.Printf("Starting server on port %d\n", port)
	err = server.ListenAndServe()
	// Wait for crashes or shutdown. Always fail first.
	if err != nil {
		app.Logger.Println(err)
	}
	app.Logger.Println("Application stopped. Bye! 👋") 

}