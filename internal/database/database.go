package database

import (
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
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

var DB *gorm.DB

func Init() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Access environment variables
	dbName := os.Getenv("DATABASE_NAME")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")

	// Print the loaded configuration for debugging
	fmt.Printf("Loaded config: DBName=%s, DBUser=%s, DBPassword=%s, DBHost=%s, DBPort=%s\n",
		dbName, dbUser, dbPassword, dbHost, dbPort)

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
		log.Fatal("Database environment variables are not set correctly")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.Person{}, &models.Course{}, &models.PersonCourse{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
}
