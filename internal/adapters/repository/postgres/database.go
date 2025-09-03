package postgres

import (
	"database/sql"
	"demoL0/internal/utils"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

// GetGormInstance создает подключение к PostgreSQL через GORM
func GetGormInstance() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем переменные окружения
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// init db
	db, err := newOrmProvider(ConnectionInfo{
		Host:     "localhost",
		Port:     utils.AtoiMust(dbPort),
		Username: dbUser,
		DBName:   dbName,
		SSLMode:  "disable",
		Password: dbPassword,
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// NewPostgresConnection создает обычное подключение к PostgreSQL
func NewPostgresConnection(info ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func newOrmProvider(info ConnectionInfo) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	s, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := s.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
