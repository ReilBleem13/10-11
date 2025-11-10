package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ReilBleem13/10-11/internal/handlers"
	"github.com/ReilBleem13/10-11/internal/repository"
	"github.com/ReilBleem13/10-11/internal/services"
	"github.com/gorilla/mux"
)

const dataDir = "data"

func main() {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("cannot create data dir:", err)
	}

	repo, err := repository.NewFileRepository()
	if err != nil {
		log.Fatal("failed to initialize file repository:", err)
	}

	srv, err := services.NewFileServices(repo)
	if err != nil {
		log.Fatal("failed to initialize services:", err)
	}

	handler, err := handlers.NewHandler(srv)
	if err != nil {
		log.Fatal("failed to initialize handler:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/check", handler.HandleCheck).Methods("POST")
	r.HandleFunc("/report", handler.HandleReport).Methods("POST")

	httpSrv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		log.Println("Server running...")
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Fatal("server error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	} else {
		log.Println("Server stopped cleanly")
	}
}
