package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/laytan/shortnr/api/link"
	"github.com/laytan/shortnr/api/user"
	"github.com/laytan/shortnr/pkg/jsonmiddleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// init runs before main
func init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
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

	linkStore := link.MysqlStorage{Conn: db}
	userStore := user.MysqlStorage{Conn: db}

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(jsonmiddleware.Middleware)

	link.SetRoutes(api, linkStore)
	user.SetRoutes(api, userStore)

	handler := cors.Default().Handler(r)

	// Start up server
	log.Fatal(http.ListenAndServe(":8080", handler))
}
