package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncorrea-13/homelab-status/internal/auth"
	"github.com/ncorrea-13/homelab-status/internal/handlers"
	"github.com/ncorrea-13/homelab-status/internal/store"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/data/status.db"
	}

	sqlStore, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	err = sqlStore.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	secret, err := auth.LoadSecret()
	if err != nil {
		log.Fatal(err)
	}

	authService := auth.NewAuth(secret)

	handler := handlers.NewHandlers(sqlStore)

	mux := handlers.NewRouter(handler, authService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		log.Println("Listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()

	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
