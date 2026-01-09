package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bikash138/students-api/internal/config"
)

func main() {
	//Load Config
	cfg := config.MustLoad()
	//Router Setup
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Server of Student"))
	})
	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}
	slog.Info("Server Started ", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func () {
		err := server.ListenAndServe() //Blocking
		if err != nil {
			log.Fatal("Server Failed to start")
		}
	} ()

	<- done

	slog.Info("Shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server stopped successfully")
}