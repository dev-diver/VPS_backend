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

	err := godotenv.Load()
	if err != nil {
		log.Printf(".env 파일을 찾을 수 없습니다. 시스템 환경 변수를 사용합니다.")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
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
