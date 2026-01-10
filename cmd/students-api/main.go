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
	"github.com/bikash138/students-api/internal/http/handlers/student"
	"github.com/bikash138/students-api/internal/storage/sqlite"
)

func main() {
	//Load Config
	cfg := config.MustLoad()

	//Db Load
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("Issue while connecting DB: %v", err)
	}
	slog.Info("DB Connected")
	//Router Setup
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))
	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server Started ", slog.String("address", cfg.Addr))

	//Gracefull Shutdown
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
	error := server.Shutdown(ctx)
	if error != nil {
		slog.Error("Failed to shutdown server", slog.String("error", error.Error()))
	}
	slog.Info("Server stopped successfully")
}