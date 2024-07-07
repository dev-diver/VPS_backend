package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func InitDB() (*Database, error) {

	err := godotenv.Load("./database/.env")
	if err != nil {
		log.Fatalf("Error loading ./databse/.env file: %v", err)
	}

	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	fmt.Println("dbUser:", dbUser, "dbPassword:", dbPassword, "dbHost:", dbHost, "dbPort:", dbPort, "dbName:", dbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	fmt.Println("db connected")
	return &Database{db}, nil
}
