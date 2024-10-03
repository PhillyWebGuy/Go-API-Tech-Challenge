package routes

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/handlers"
)

// RegisterRoutes registers the routes for the application.
func RegisterRoutes(r *chi.Mux) {
	r.Get("/api/person", handlers.GetPersons)
	r.Post("/api/person", handlers.CreatePerson)
	r.Get("/api/person/{name}", handlers.GetPerson)
	r.Delete("/api/person/{name}", handlers.DeletePerson)
	r.Put("/api/person/{name}", handlers.UpdatePerson)

	r.Get("/api/course", handlers.GetCourses)
	r.Get("/api/course/{id}", handlers.GetCourse)
	r.Put("/api/course/{id}", handlers.UpdateCourse)
	r.Post("/api/course", handlers.CreateCourse)
	r.Delete("/api/course/{id}", handlers.DeleteCourse)

	routes := r.Routes()
	for _, route := range routes {
		fmt.Printf("Registered route: %s\n", route.Pattern)
	}
}
