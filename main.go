package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/laytan/shortnr/pkg/jsonmiddleware"

	"github.com/laytan/shortnr/internal/link"
	"github.com/laytan/shortnr/internal/redirect"
	"github.com/laytan/shortnr/internal/user/routes"
	userStorage "github.com/laytan/shortnr/internal/user/storage"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	linkStorage "github.com/laytan/shortnr/internal/storage"
)

// init runs before main
func init() {
	// Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
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

	r := mux.NewRouter()

	linksStore := linkStorage.MysqlStorage{Conn: db}
	usersStore := userStorage.Mysql{Conn: db}

	// General routes
	redirect.SetRoutes(r, linksStore)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(jsonmiddleware.Middleware)

	// API routes
	link.SetRoutes(api, linksStore)
	routes.Set(api, usersStore)

	// Start up server
	log.Fatal(http.ListenAndServe(":8080", r))
}
