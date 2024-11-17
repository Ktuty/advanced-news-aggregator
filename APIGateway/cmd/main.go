package main

import (
	"api-getaway/internal/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	handler := handlers.NewHandler()
	go func() {
		if err := http.ListenAndServe(":80", handler.InitRouts()); err != nil {
			log.Fatalf("error starting http server: %s", err.Error())
		}
	}()

	log.Println("Start running server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
}
