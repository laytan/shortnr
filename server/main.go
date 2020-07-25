package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

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
	r.Use(jsonmiddleware.Middleware)

	usersRouter := r.PathPrefix("/api/v1/users").Subrouter()
	linksRouter := r.PathPrefix("/api/v1/links").Subrouter()

	linksStore := link.MysqlStorage{Conn: db}
	usersStore := user.MysqlStorage{Conn: db}

	link.SetRoutes(linksRouter, linksStore)
	user.SetRoutes(usersRouter, usersStore)

	// apiRouter := r.PathPrefix("/api/v1").Subrouter()
	// apiRouter.Use(jsonmiddleware.Middleware)

	// link.SetRoutes(apiRouter, linkStore)
	// user.SetRoutes(apiRouter, userStore)

	// apiAuthRouter := r.PathPrefix("").Subrouter()
	// apiAuthRouter.Use(user.JwtAuthorization)
	// user.SetAuthRoutes(apiAuthRouter, userStore)

	handler := cors.Default().Handler(r)

	// Start up server
	log.Fatal(http.ListenAndServe(":8080", handler))
}
