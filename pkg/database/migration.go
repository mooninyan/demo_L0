package database

import (
	"database/sql"
	"demoL0/internal/utils"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"strings"
)

func RunMigrations() {
	migrationPath := "file://./pkg/database/migrations"

	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dName := os.Getenv("DB_NAME")
	databaseURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPassword, dbPort, dName)

	m, err := migrate.New(migrationPath, databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}

func InitialMigration() {
	db, _ := defaultConnection()
	defer db.Close()
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	migration := "CREATE USER $2 WITH PASSWORD '$3';"
	migration = strings.ReplaceAll(migration, "$3", dbPassword)
	migration = strings.ReplaceAll(migration, "$2", dbUser)
	_, _ = db.Exec(migration)

	migrationInit := "CREATE DATABASE " + dbName + " WITH OWNER " + dbUser
	_, _ = db.Exec(migrationInit)

	log.Println("Initial db created/exists")
}

func defaultConnection() (*sql.DB, error) {
	defaultDbPort := os.Getenv("DEFAULT_DB_PORT")
	defaultDbUser := os.Getenv("DEFAULT_DB_USER")
	defaultDbPassword := os.Getenv("DEFAULT_DB_PASSWORD")
	defaultDbName := os.Getenv("DEFAULT_DB_NAME")

	db, err := NewPostgresConnection(ConnectionInfo{
		Host:     "localhost",
		Port:     utils.AtoiMust(defaultDbPort),
		Username: defaultDbUser,
		DBName:   defaultDbName,
		SSLMode:  "disable",
		Password: defaultDbPassword,
	})
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
