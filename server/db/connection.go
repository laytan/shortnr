package db

import (
	"database/sql"
	"fmt"
	"os"

	// Mysql driver has to be imported
	_ "github.com/go-sql-driver/mysql"
)

// GetConnection establishes a db connection and returns it
func GetConnection() *sql.DB {
	// Get connection vars out of env
	username, exists := os.LookupEnv("DB_USERNAME")
	if !exists {
		panic("No DB_USERNAME env variable set")
	}
	password, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		panic("No DB_PASSWORD env variable set")
	}
	database, exists := os.LookupEnv("DB_DATABASE")
	if !exists {
		panic("No DB_DATABASE env variable set")
	}

	// Connect to mysql
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", username, password, database))
	if err != nil {
		panic(err.Error())
	}

	return db
}
