package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
)

func setupTestCourseDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Course{}, &models.Person{})
	return db
}

func TestGetCourses(t *testing.T) {
	db := setupTestCourseDB()
	handler := NewRequestHandler(db)

	// Insert test data
	db.Create(&models.Course{Name: "Course 1"})
	db.Create(&models.Course{Name: "Course 2"})

	req, _ := http.NewRequest("GET", "/courses", nil)
	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.GetCourses)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var courses []models.Course
	json.Unmarshal(rr.Body.Bytes(), &courses)
	assert.Equal(t, 2, len(courses))
}

func TestGetCourse(t *testing.T) {
	db := setupTestCourseDB()
	handler := NewRequestHandler(db)

	// Insert test data
	course := models.Course{Name: "Course 1"}
	db.Create(&course)

	req, _ := http.NewRequest("GET", "/courses/"+strconv.Itoa(int(course.ID)), nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/courses/{id}", handler.GetCourse)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedCourse models.Course
	json.Unmarshal(rr.Body.Bytes(), &retrievedCourse)
	assert.Equal(t, course.Name, retrievedCourse.Name)
}

func TestCreateCourse(t *testing.T) {
	db := setupTestCourseDB()
	handler := NewRequestHandler(db)

	course := models.Course{Name: "Course 1"}
	body, _ := json.Marshal(course)
	req, _ := http.NewRequest("POST", "/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.CreateCourse)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var createdCourse models.Course
	json.Unmarshal(rr.Body.Bytes(), &createdCourse)
	assert.Equal(t, course.Name, createdCourse.Name)
}

func TestUpdateCourse(t *testing.T) {
	db := setupTestCourseDB()
	handler := NewRequestHandler(db)

	// Insert test data
	course := models.Course{Name: "Course 1"}
	db.Create(&course)

	updatedCourse := models.Course{Name: "Updated Course"}
	body, _ := json.Marshal(updatedCourse)
	req, _ := http.NewRequest("PUT", "/courses/"+strconv.Itoa(int(course.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Put("/courses/{id}", handler.UpdateCourse)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var retrievedCourse models.Course
	db.First(&retrievedCourse, course.ID)
	assert.Equal(t, updatedCourse.Name, retrievedCourse.Name)
}

func TestDeleteCourse(t *testing.T) {
	db := setupTestCourseDB()
	handler := NewRequestHandler(db)

	// Insert test data
	course := models.Course{Name: "Course 1"}
	db.Create(&course)

	// Insert associated records in person_course table
	person := models.Person{FirstName: "John", LastName: "Doe"}
	db.Create(&person)
	db.Exec("INSERT INTO person_course (person_id, course_id) VALUES (?, ?)", person.ID, course.ID)

	// Check that the record exists in person_course table
	var personCourseCount int64
	db.Table("person_course").Where("course_id = ?", course.ID).Count(&personCourseCount)
	assert.Equal(t, int64(1), personCourseCount)

	req, _ := http.NewRequest("DELETE", "/courses/"+strconv.Itoa(int(course.ID)), nil)
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Delete("/courses/{id}", handler.DeleteCourse)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Check that the course is deleted
	var retrievedCourse models.Course
	err := db.First(&retrievedCourse, course.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	// Check that the associated records in person_course table are deleted
	db.Table("person_course").Where("course_id = ?", course.ID).Count(&personCourseCount)
	assert.Equal(t, int64(0), personCourseCount)
}
