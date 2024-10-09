package routes

import (
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/handlers"
)

// RegisterRoutes registers the routes for the application.
func RegisterRoutes(r *chi.Mux, requestHandler *handlers.RequestHandler) {
	r.Get("/api/person", requestHandler.GetPersons)
	r.Post("/api/person", requestHandler.CreatePerson)
	r.Get("/api/person/{name}", requestHandler.GetPerson)
	r.Delete("/api/person/{name}", requestHandler.DeletePerson)
	r.Put("/api/person/{name}", requestHandler.UpdatePerson)

	r.Get("/api/course", requestHandler.GetCourses)
	r.Get("/api/course/{id}", requestHandler.GetCourse)
	r.Put("/api/course/{id}", requestHandler.UpdateCourse)
	r.Post("/api/course", requestHandler.CreateCourse)
	r.Delete("/api/course/{id}", requestHandler.DeleteCourse)

	routes := r.Routes()
	for _, route := range routes {
		log.Printf("Registered route: %s\n", route.Pattern)
	}
}
