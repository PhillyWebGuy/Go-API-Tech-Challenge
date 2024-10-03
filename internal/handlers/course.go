package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/database"
	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
)

// GetCourses handles the request to get all courses.
// It retrieves all courses from the database and returns them as a JSON response.
//
// @param w http.ResponseWriter - the response writer to send the response
// @param r *http.Request - the request to get all courses
//
// @response 200 - Courses retrieved successfully
// @response 500 - Internal server error
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

// GetCourse handles the request to get a course by ID.
// It retrieves the course ID from the URL, converts it to an integer,
// finds the course in the database, and returns the course as a JSON response.
//
// @param w http.ResponseWriter - the response writer to send the response
// @param r *http.Request - the request containing the course ID in the URL
//
// @response 200 - Course retrieved successfully
// @response 400 - Invalid course ID
// @response 404 - Course not found
func GetCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Convert the ID to an integer
	intID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find the course by ID
	var course models.Course
	if err := database.DB.First(&course, "id = ?", intID).Error; err != nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

// UpdateCourse handles the request to update a specific course.
// It decodes the request body into a Course struct, finds the existing course by ID,
// updates the course fields, and saves the updated course to the database.
//
// @param w http.ResponseWriter - the response writer to send the response
// @param r *http.Request - the request containing the course data and ID in the URL
//
// @response 200 - Course updated successfully
// @response 400 - Invalid request payload
// @response 404 - Course not found
// @response 500 - Internal server error
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
// It starts a new transaction, finds the course by ID, deletes associated records
// in the person_course table, deletes the course, and commits the transaction.
//
// @param w http.ResponseWriter - the response writer to send the response
// @param r *http.Request - the request containing the course ID in the URL
//
// @response 200 - Course deleted successfully
// @response 404 - Course not found
// @response 500 - Internal server error
func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Start a new transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Find the course by ID
	var course models.Course
	if err := tx.First(&course, "id = ?", id).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	// Delete records in person_course that are associated with the course
	if err := tx.Exec("DELETE FROM person_course WHERE course_id = ?", id).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the course
	if err := tx.Delete(&course).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Course deleted successfully"})
}

// CreateCourse handles the request to create a new course.
// It decodes the request body into a Course struct, saves it to the database,
// and returns the created course in the response.
//
// @param w http.ResponseWriter - the response writer to send the response
// @param r *http.Request - the request containing the course data
//
// @response 201 - Course created successfully
// @response 400 - Invalid request payload
// @response 500 - Internal server error
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	// Decode the request body into the course struct
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Start a new transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Check if the course already exists in the database
	var existingCourse models.Course
	if err := tx.Where("name = ?", course.Name).First(&existingCourse).Error; err == nil {
		tx.Rollback()
		http.Error(w, "Course already exists", http.StatusConflict)
		return
	}

	// Save the course to the database
	if err := tx.Create(&course).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}
