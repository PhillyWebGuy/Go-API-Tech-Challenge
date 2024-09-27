package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/database"
	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	// Initialize a new in-memory SQLite database for testing
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.DB = db

	// Migrate the schema
	database.DB.AutoMigrate(&models.Person{})
}

func TestGetPersons(t *testing.T) {
	setupTestDB()

	// Insert test data
	database.DB.Create(&models.Person{FirstName: "John", LastName: "Doe", Type: "student", Age: 20})
	database.DB.Create(&models.Person{FirstName: "Jane", LastName: "Smith", Type: "professor", Age: 45})

	req, err := http.NewRequest("GET", "/api/person", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPersons)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var persons []models.Person
	err = json.NewDecoder(rr.Body).Decode(&persons)
	assert.NoError(t, err)
	assert.Len(t, persons, 2)
}

func TestCreatePerson(t *testing.T) {
	setupTestDB()

	person := models.Person{FirstName: "Alice", LastName: "Johnson", Type: "student", Age: 22}
	body, err := json.Marshal(person)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/person", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePerson)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdPerson models.Person
	err = json.NewDecoder(rr.Body).Decode(&createdPerson)
	assert.NoError(t, err)
	assert.Equal(t, person.FirstName, createdPerson.FirstName)
	assert.Equal(t, person.LastName, createdPerson.LastName)
	assert.Equal(t, person.Type, createdPerson.Type)
	assert.Equal(t, person.Age, createdPerson.Age)
}
