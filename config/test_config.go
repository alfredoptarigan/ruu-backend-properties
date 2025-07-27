package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"alfredo/ruu-properties/pkg/models"
)

func InitTestDatabase() *gorm.DB {
	// Force read environment variables directly
	testDBHost := "localhost"
	testDBPort := "5432"
	testDBUser := "alfredopatriciustarigan"
	testDBPassword := "test" // Set default password
	testDBName := "test_ruu_properties"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		testDBHost, testDBUser, testDBPassword, testDBName, testDBPort)

	// Debug: print DSN to verify
	fmt.Printf("Test DB DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent mode for tests
	})

	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Enable UUID extension for PostgreSQL
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Auto migrate for tests
	err = db.AutoMigrate(&models.User{}, &models.Client{}) // Add all your models here
	if err != nil {
		log.Fatal("Failed to migrate test database:", err)
	}

	// Verify table creation
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'users'").Scan(&count)
	if count == 0 {
		log.Fatal("Users table was not created successfully")
	}

	fmt.Println("Test database initialized successfully")
	return db
}
