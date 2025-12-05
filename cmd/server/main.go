package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	config "github.com/user/todo_auth"
	"github.com/user/todo_auth/internal/database"
	"github.com/user/todo_auth/internal/handlers"
)

func main() {
	// Загружаем конфигурацию из .env.development и переменных окружения
	appCfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	dbCfg := database.Config{
		DSN: appCfg.DatabaseDSN,
	}

	db, err := database.NewConnection(dbCfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.Migrate(context.Background(), db); err != nil {
		slog.Error("database migration failed", "error", err)
		os.Exit(1)
	}

	router := mux.NewRouter()

	handlers.RegisterRoutes(router, db)

	slog.Info(fmt.Sprintf("Starting server at http://localhost%s", appCfg.ServerPort))
	if err := http.ListenAndServe(appCfg.ServerPort, router); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
