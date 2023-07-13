package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jesiahharris/rss-agg/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("portString does not exist")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL does not exist")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to db: ", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// create router and check status with handler
	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	// mount and set path for router. path will be /v1/*
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v \n", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Port: %s \n", dbURL)
}
