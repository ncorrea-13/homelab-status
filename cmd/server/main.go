package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ncorrea-13/homelab-status/internal/auth"
	"github.com/ncorrea-13/homelab-status/internal/handlers"
	"github.com/ncorrea-13/homelab-status/internal/store"
)

func main() {
	ctx := context.Background()

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

	log.Println("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
