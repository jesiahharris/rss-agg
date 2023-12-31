package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jesiahharris/rss-agg/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

const htmlIndex = `<html><head>
<title>
Test input form
</title>
</head>
<body> 
 <p> Input API Key: </p><input type ="text" placeholder="API Key" id="key" />
</html>
`

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
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

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	// call startScraping() on new go routine before listenAndServe
	go startScraping(db, 10, time.Minute)

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

	v1Router.Route("/users", func(r chi.Router) {
		r.Get("/", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
		r.Post("/", apiCfg.handlerCreateUser)
	})

	v1Router.Route("/feeds", func(r chi.Router) {
		r.Post("/", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
		r.Get("/", apiCfg.handlerGetFeeds)
		r.Delete("/{feedID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeed))
	})

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	v1Router.Route("/feed_follows", func(r chi.Router) {
		r.Post("/", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
		r.Get("/", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
		r.Delete("/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	})

	v1Router.HandleFunc("/", handleIndex)

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
