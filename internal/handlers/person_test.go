package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/database"
	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Person{}, &models.PersonCourse{}, &models.Course{})
	return db
}

func TestGetPersons(t *testing.T) {
	db := setupTestDB()
	database.DB = db

	// Insert test data
	db.Create(&models.Person{FirstName: "John", LastName: "Doe", Type: "student", Age: 25})
	db.Create(&models.Person{FirstName: "Jane", LastName: "Smith", Type: "student", Age: 35})

	req, _ := http.NewRequest("GET", "/persons", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPersons)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var persons []models.Person
	json.Unmarshal(rr.Body.Bytes(), &persons)
	assert.Equal(t, 2, len(persons))
}

func TestGetPerson(t *testing.T) {
	db := setupTestDB()
	database.DB = db

	// Insert test data
	person := models.Person{FirstName: "John", LastName: "Doe", Type: "student", Age: 25}
	db.Create(&person)

	req, _ := http.NewRequest("GET", "/persons/"+url.QueryEscape("John Doe"), nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/persons/{name}", GetPerson)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedPerson models.Person
	json.Unmarshal(rr.Body.Bytes(), &retrievedPerson)
	assert.Equal(t, person.FirstName, retrievedPerson.FirstName)
	assert.Equal(t, person.LastName, retrievedPerson.LastName)
}

func TestCreatePerson(t *testing.T) {
	db := setupTestDB()
	database.DB = db

	person := models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
	}

	personWithCourses := models.PersonWithCourses{
		Person:  person,
		Courses: []int{1, 2},
	}

	body, _ := json.Marshal(personWithCourses)
	req, _ := http.NewRequest("POST", "/persons", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePerson)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdPerson models.Person
	json.Unmarshal(rr.Body.Bytes(), &createdPerson)
	assert.Equal(t, person.FirstName, createdPerson.FirstName)
	assert.Equal(t, person.LastName, createdPerson.LastName)

	// Check that the records were inserted into the person_course table
	var personCourses []models.PersonCourse
	database.DB.Where("person_id = ?", createdPerson.ID).Find(&personCourses)
	assert.Equal(t, 2, len(personCourses))
}

func TestDeletePerson(t *testing.T) {
	db := setupTestDB()
	database.DB = db

	// Insert test data
	person := models.Person{FirstName: "John", LastName: "Doe", Type: "student", Age: 25}
	db.Create(&person)

	// Insert associated records in person_course table
	course := models.Course{Name: "Course 1"}
	db.Create(&course)
	db.Exec("INSERT INTO person_course (person_id, course_id) VALUES (?, ?)", person.ID, course.ID)

	req, _ := http.NewRequest("DELETE", "/persons/"+url.QueryEscape("John Doe"), nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Delete("/persons/{name}", DeletePerson)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the person is deleted
	var retrievedPerson models.Person
	err := db.First(&retrievedPerson, person.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Check that the associated records in person_course table are deleted
	var personCourses []models.PersonCourse
	database.DB.Where("person_id = ?", person.ID).Find(&personCourses)
	assert.Equal(t, 0, len(personCourses))
}

func TestUpdatePerson(t *testing.T) {
	db := setupTestDB()
	database.DB = db

	// Insert test data
	person := models.Person{
		FirstName: "Steve",
		LastName:  "Jones",
		Type:      "student",
		Age:       30,
	}

	db.Create(&person)

	// Insert associated records in person_course table
	course := models.Course{Name: "Course 1"}
	db.Create(&course)
	db.Exec("INSERT INTO person_course (person_id, course_id) VALUES (?, ?)", person.ID, course.ID)

	/*newUpdatedPerson := models.Person{
		FirstName: "Steve",
		LastName:  "Jones",
		Type:      "student",
		Age:       30,
	}*/

	updatedPersonWithCourses := models.PersonWithCourses{
		Person:  person,
		Courses: []int{1, 2},
	}
	body, _ := json.Marshal(updatedPersonWithCourses)
	req, _ := http.NewRequest("PUT", "/persons/"+url.QueryEscape("Steve Jones"), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Put("/persons/{name}", UpdatePerson)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var updatedPerson models.Person
	database.DB.First(&updatedPerson, person.ID)
	assert.Equal(t, updatedPersonWithCourses.FirstName, updatedPerson.FirstName)
	assert.Equal(t, updatedPersonWithCourses.LastName, updatedPerson.LastName)
	assert.Equal(t, updatedPersonWithCourses.Type, updatedPerson.Type)
	assert.Equal(t, updatedPersonWithCourses.Age, updatedPerson.Age)

	// Check that the records were updated in the person_course table
	var personCourses []models.PersonCourse
	database.DB.Where("person_id = ?", updatedPerson.ID).Find(&personCourses)
	assert.Equal(t, 1, len(personCourses))
}
