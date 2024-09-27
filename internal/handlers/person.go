package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/database"
	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// GetPersons handles the request to get a list of persons.
func GetPersons(w http.ResponseWriter, r *http.Request) {
	var persons []models.Person

	// Query the database for all persons
	if err := database.DB.Find(&persons).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

// GetPerson handles the request to retrieve a person by ID.
func GetPerson(w http.ResponseWriter, r *http.Request) {
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
	if err := database.DB.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&person).Error; err != nil {
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
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person

	// Decode the request body into the person struct
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert the person into the database
	if err := database.DB.Create(&person).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Debug logging
	log.Printf("Status Code: %d", http.StatusCreated)
	log.Printf("Response: %+v", person)

	json.NewEncoder(w).Encode(person)
}

// DeletePerson handles the request to delete a person by first name and last name.
func DeletePerson(w http.ResponseWriter, r *http.Request) {
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
	if err := database.DB.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Person not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Delete the person from the database
	if err := database.DB.Delete(&person).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Debug logging
	log.Printf("Status Code: %d", http.StatusOK)
	log.Printf("Deleted Person: %+v", person)

	json.NewEncoder(w).Encode(map[string]string{"message": "Person deleted successfully"})
}

// UpdatePerson handles the request to update a person's details.
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
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
	if err := database.DB.Where("first_name = ? AND last_name = ?", firstName, lastName).First(&person).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Person not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Decode the request body into the person struct
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Save the updated person details to the database
	if err := database.DB.Save(&person).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}
