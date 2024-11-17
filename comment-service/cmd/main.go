package main

import (
	"comment-service/internal/handlers"
	"comment-service/internal/repository"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("error connecting to database: %s", err.Error())
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	handler := handlers.NewHandler(repo)

	go func() {
		if err := http.ListenAndServe(":8082", handler.InitRouts()); err != nil {
			log.Fatalf("error starting http server: %s", err.Error())
		}
	}()

	log.Println("Start running server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
}
