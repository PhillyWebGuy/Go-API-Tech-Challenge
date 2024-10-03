package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/routes"
)

func NewServer() {
	r := SetupRouter()

	port := ":8000"
	log.Printf("Starting server on %s\n", port)

	if err := http.ListenAndServe(port, r); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
}

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	routes.RegisterRoutes(r)

	return r
}
