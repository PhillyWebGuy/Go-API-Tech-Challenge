package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PhillyWebGuy/Go-API-Tech-Challenge/internal/models"
)

func LoadConfig() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	return nil
}

func Init() (*gorm.DB, error) {
	err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Access environment variables
	dbName := os.Getenv("DATABASE_NAME")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")

	// Print the loaded configuration for debugging
	log.Printf("Loaded config: DBName=%s, DBUser=%s, DBPassword=%s, DBHost=%s, DBPort=%s\n",
		dbName, dbUser, dbPassword, dbHost, dbPort)

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		return nil, errors.New("database environment variables are not set correctly")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.Person{}, &models.Course{}, &models.PersonCourse{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	return db, nil
}
