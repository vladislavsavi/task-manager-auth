package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	config "github.com/user/todo_auth"
	"github.com/user/todo_auth/internal/database"
	"github.com/user/todo_auth/internal/handlers"
)

func main() {
	// Загружаем конфигурацию из .env.development и переменных окружения
	appCfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	dbCfg := database.Config{
		DSN: appCfg.DatabaseDSN,
	}

	db, err := database.NewConnection(dbCfg) //nolint:all
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(context.Background(), db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	router := mux.NewRouter()

	handlers.RegisterRoutes(router, db)

	log.Printf("Starting server on port %s", appCfg.ServerPort)
	log.Fatal(http.ListenAndServe(appCfg.ServerPort, router))
}
