package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gopalpokhrel1/students-api/internal/config"
	student "github.com/gopalpokhrel1/students-api/internal/http/handlers/students"
	"github.com/gopalpokhrel1/students-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//setup database

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

	//setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}
	fmt.Println("Server started at localhost:8082")

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			fmt.Println("Error occured")
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	slog.Info("Shutting down the server")

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown server")
	}

	slog.Info("Server shutdown successfully")
}
