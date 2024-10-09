package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
)

func NewTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err := db.AutoMigrate(&models.Person{}, &models.PersonCourse{}, &models.Course{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return db
}
