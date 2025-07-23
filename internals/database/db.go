package database

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Client() *gorm.DB {
	return DB
}
func Connect() {
	_ = godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s database=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	
	// Implement retry logic for database connection
	var db *gorm.DB
	var err error
	for retries := 5; retries > 0; retries-- {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("Unable to connect to database. Retrying in 5 seconds... (%d attempts left)\n", retries)
		time.Sleep(5 * time.Second)
	}
	
	if err != nil {
		fmt.Printf("Failed to connect to database after multiple attempts: %v\n", err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Unable to get sql database from gorm: %v\n", err)
		panic(err)
	}
	
	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("Unable to ping database: %v\n", err)
		panic(err)
	}

	DB = db
	fmt.Println("Connected to database successfully")
	fmt.Println("Running Migrations...")

}
