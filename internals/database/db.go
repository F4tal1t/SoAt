package database

import (
	"fmt"
	"os"

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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Unable to connect to database. DSN: %s\nError: %v\n", dsn, err)
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Unable to get sql database from gorm , error : %v", err)
		panic(err)
	}
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("Unable to ping database, error : %v", err)
		panic(err)
	}

	DB = db
	fmt.Println("Connected to database successfully")
	fmt.Println("Running Migrations...")

}
