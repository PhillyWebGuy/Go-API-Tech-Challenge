package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/database"
	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
)

// GetCourses handles the request to get all courses
func GetCourses(w http.ResponseWriter, r *http.Request) {
	var courses []models.Course

	// Retrieve all courses without preloading the associated people
	result := database.DB.Find(&courses)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// GetCourse handles the request to get a course by ID and includes people in the course
func GetCourse(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL parameters
	id := chi.URLParam(r, "id")

	// Convert the ID to an unsigned integer
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		http.Error(w, "Invalid course ID", http.StatusBadRequest)
		return
	}

	// Cast the parsed ID to uint
	courseID := uint(parsedID)

	var course models.Course

	// Retrieve the course from the database
	result := database.DB.First(&course, "id = ?", courseID)
	if result.Error != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	var people []models.Person

	// Retrieve the people associated with the course
	result = database.DB.Table("people").Select("people.*").
		Joins("JOIN person_course ON person_course.person_id = people.id").
		Where("person_course.course_id = ?", courseID).Find(&people)
	if result.Error != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Create a response struct to include course and people
	response := struct {
		Course models.Course   `json:"course"`
		People []models.Person `json:"people"`
	}{
		Course: course,
		People: people,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateCourse handles the request to update a specific course.
func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var course models.Course

	// Decode the request body into the course struct
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the course by ID
	var existingCourse models.Course
	if err := database.DB.First(&existingCourse, "id = ?", id).Error; err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	// Update the course fields
	existingCourse.Name = course.Name

	// Save the updated course to the database
	if err := database.DB.Save(&existingCourse).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingCourse)
}

// DeleteCourse handles the request to delete a specific course.
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Find the course by ID
	var course models.Course
	if err := database.DB.First(&course, "id = ?", id).Error; err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	// Delete the course
	if err := database.DB.Delete(&course).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Course deleted successfully"})
}

// CreateCourse handles the request to create a new course.
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	// Decode the request body into the course struct
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Save the course to the database
	if err := database.DB.Create(&course).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}
