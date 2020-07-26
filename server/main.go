package main

import (
	"log"
	"net/http"

	"github.com/laytan/shortnr/db"

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
	db := db.GetConnection()

	r := mux.NewRouter()
	r.Use(jsonmiddleware.Middleware)

	usersRouter := r.PathPrefix("/api/v1/users").Subrouter()
	linksRouter := r.PathPrefix("/api/v1/links").Subrouter()

	linksStore := link.MysqlStorage{Conn: db}
	usersStore := user.MysqlStorage{Conn: db}

	link.SetRoutes(linksRouter, linksStore)
	user.SetRoutes(usersRouter, usersStore)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "DELETE", "POST"},
	})

	// Insert the middleware
	handler := c.Handler(r)

	// Start up server
	log.Fatal(http.ListenAndServe(":8080", handler))
}
