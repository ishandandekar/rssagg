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
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("hello world")

	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	fmt.Println(portString)

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	fmt.Println(dbURL)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	queries := database.New(conn)
	apicfg := apiConfig{DB: queries}

	go startScraping(queries, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			},
		),
	)
	srv := &http.Server{Handler: router, Addr: ":" + portString}

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Post("/users", apicfg.handlerCreateUser)
	v1Router.Get("/users", apicfg.middlewareAuth(apicfg.handlerGetUser))

	v1Router.Post("/feeds", apicfg.middlewareAuth(apicfg.handlerCreateFeed))
	v1Router.Get("/feeds", apicfg.handlerGetFeeds)

	v1Router.Post("/posts", apicfg.middlewareAuth(apicfg.handlerGetPostsForUser))

	v1Router.Post("/feed_follows", apicfg.middlewareAuth(apicfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apicfg.middlewareAuth(apicfg.handlerGetFeedFollows))
	v1Router.Delete(
		"/feed_follows/{feedFollowID}",
		apicfg.middlewareAuth(apicfg.handlerDeleteFeedFollow),
	)

	router.Mount("/v1", v1Router)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
