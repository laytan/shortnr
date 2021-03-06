package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/laytan/shortnr/db"

	_ "github.com/go-sql-driver/mysql"

	"github.com/laytan/shortnr/api/click"
	"github.com/laytan/shortnr/api/link"
	"github.com/laytan/shortnr/api/user"
	"github.com/laytan/shortnr/pkg/jsonmiddleware"
	"github.com/laytan/shortnr/pkg/responder"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// init runs before main
func init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, is this production?")
	}
}

func main() {
	db := db.GetConnection()

	r := mux.NewRouter()
	r.Use(jsonmiddleware.Middleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		responder.Res{
			Message: fmt.Sprintf("Hello, world! You should use our app at: %s", os.Getenv("FRONT_END_URL")),
		}.Send(w)
	})

	usersRouter := r.PathPrefix("/api/v1/users").Subrouter()
	linksRouter := r.PathPrefix("/api/v1/links").Subrouter()
	clicksRouter := r.PathPrefix("/api/v1/clicks").Subrouter()

	linksStore := link.MysqlStorage{Conn: db}
	usersStore := user.MysqlStorage{Conn: db}
	clicksStore := click.MysqlStorage{Conn: db}

	link.SetRoutes(linksRouter, linksStore, clicksStore)
	user.SetRoutes(usersRouter, usersStore)
	click.SetRoutes(clicksRouter, clicksStore)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONT_END_URL")},
		AllowCredentials: true,
		AllowedMethods:   []string{"OPTIONS", "HEAD", "GET", "DELETE", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	// Insert the middleware
	handler := c.Handler(r)

	// Start up server
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
}
