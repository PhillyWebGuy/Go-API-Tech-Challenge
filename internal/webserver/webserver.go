package webserver

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/handlers"
	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/routes"
)

// NewServer initializes a new web server with the given database connection.
func NewServer(db *gorm.DB) {
	// Initialize your routes and handlers here, passing the db connection as needed
	requestHandler := handlers.NewRequestHandler(db)
	r := SetupRouter(requestHandler)

	const port = ":8000"
	log.Printf("Starting server on %s\n", port)

	const readTimeout = 5 * time.Second
	const writeTimeout = 10 * time.Second
	const idleTimeout = 15 * time.Second

	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
}

func SetupRouter(requestHandler *handlers.RequestHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	routes.RegisterRoutes(r, requestHandler)

	return r
}
