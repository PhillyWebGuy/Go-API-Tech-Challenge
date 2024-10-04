package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
)

func validatePersonWithCourses(personWithCourses models.PersonWithCourses) error {
	return models.Validate.Struct(personWithCourses)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}

// GetPersons handles the request to get a list of persons.
func (h *RequestHandler) GetPersons(w http.ResponseWriter, r *http.Request) {
	var persons []models.Person

	// Query the database for all persons
	if err := h.DB.Find(&persons).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

// GetPerson handles the request to retrieve a person by ID.
func (h *RequestHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	// Extract the name from the URL parameters
	name := chi.URLParam(r, "name")

	// Decode the URL-encoded name
	decodedFullName, err := url.QueryUnescape(name)
	if err != nil {
		http.Error(w, "Invalid full name format", http.StatusBadRequest)
		return
	}

	// Split the full name into first name and last name
	names := strings.SplitN(decodedFullName, " ", 2)
	if len(names) != 2 {
		http.Error(w, "Invalid full name format", http.StatusBadRequest)
		return
	}
	firstName := names[0]
	lastName := names[1]

	var person models.Person

	// Retrieve the person from the database
	if err := h.DB.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Person not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// CreatePerson handles the request to create a new person.
func (h *RequestHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var personWithCourses models.PersonWithCourses

	// Decode the request body into the personWithCourses struct
	if err := json.NewDecoder(r.Body).Decode(&personWithCourses); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validatePersonWithCourses(personWithCourses); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Start a new transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Check if the person already exists in the database
	var existingPerson models.Person
	if err := tx.Where("first_name = ? AND last_name = ?", personWithCourses.FirstName, personWithCourses.LastName).First(&existingPerson).Error; err == nil {
		tx.Rollback()
		http.Error(w, "Person already exists", http.StatusConflict)
		return
	}

	// Create a new Person struct from the personWithCourses struct
	person := models.Person{
		FirstName: personWithCourses.FirstName,
		LastName:  personWithCourses.LastName,
		Type:      personWithCourses.Type,
		Age:       personWithCourses.Age,
	}

	// Save the person to the database
	if err := tx.Create(&person).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert records into the person_course table
	for _, courseID := range personWithCourses.Courses {
		if err := tx.Exec("INSERT INTO person_course (person_id, course_id) VALUES (?, ?)", person.ID, courseID).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

// DeletePerson handles the request to delete a person by first name and last name.
func (h *RequestHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	// Extract the name from the URL parameters
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}

	// Decode the URL-encoded full name
	decodedFullName, err := url.QueryUnescape(name)
	if err != nil {
		http.Error(w, "Invalid full name format", http.StatusBadRequest)
		return
	}

	// Split the full name into first name and last name
	names := strings.SplitN(decodedFullName, " ", 2)
	if len(names) != 2 {
		http.Error(w, "Invalid full name format", http.StatusBadRequest)
		return
	}
	firstName := names[0]
	lastName := names[1]

	var person models.Person

	// Start a new transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Retrieve the person from the database within the transaction
	if err := tx.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			http.Error(w, "Person not found", http.StatusNotFound)
		} else {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Explicitly delete related person_course records within the transaction
	if err := tx.Where("person_id = ?", person.ID).Delete(&models.PersonCourse{}).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the person from the database within the transaction
	if err := tx.Delete(&person).Error; err != nil {
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
	json.NewEncoder(w).Encode(map[string]string{"message": "Person deleted successfully"})
}

// UpdatePerson handles the request to update a person's details.
func (h *RequestHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	// Extract the name from the URL parameters
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Name parameter is missing", http.StatusBadRequest)
		return
	}

	// Decode the URL-encoded full name
	decodedFullName, err := url.QueryUnescape(name)
	if err != nil {
		http.Error(w, "Invalid full name format", http.StatusBadRequest)
		return
	}

	// Split the full name into first name and last name
	names := strings.SplitN(decodedFullName, " ", 2)
	if len(names) != 2 {
		http.Error(w, "Invalid full name format", http.StatusBadRequest)
		return
	}
	firstName := names[0]
	lastName := names[1]

	var person models.Person

	// Retrieve the person from the database
	if err := h.DB.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Person not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var personWithCourses models.PersonWithCourses

	// Decode the request body into the personWithCourses struct
	if err := json.NewDecoder(r.Body).Decode(&personWithCourses); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validatePersonWithCourses(personWithCourses); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Start a new transaction
	tx := h.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Update the person's details
	/*person.FirstName = personWithCourses.FirstName
	  person.LastName = personWithCourses.LastName
	  person.Type = personWithCourses.Type
	  person.Age = personWithCourses.Age*/

	// Save the updated person details to the database
	if err := tx.Save(&personWithCourses).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete existing records in the person_course table for the person
	if err := tx.Exec("DELETE FROM person_course WHERE person_id = ?", person.ID).Error; err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert new records into the person_course table
	for _, courseID := range personWithCourses.Courses {
		if err := tx.Exec("INSERT INTO person_course (person_id, course_id) VALUES (?, ?)", person.ID, courseID).Error; err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}
