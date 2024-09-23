package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var courses = map[string]Course{
	"1": {ID: "1", Name: "Go Programming"},
	"2": {ID: "2", Name: "Web Development"},
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/course", getCourses)
	r.Get("/api/course/{id}", getCourse)
	r.Put("/api/course/{id}", updateCourse)

	http.ListenAndServe(":8000", r)
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	course, ok := courses[id]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func updateCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	course.ID = id
	courses[id] = course
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}
