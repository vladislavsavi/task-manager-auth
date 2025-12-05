package handlers

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type Handler struct {
	DB *sql.DB
}

// RegisterRoutes - регистрирует все маршруты
func RegisterRoutes(router *mux.Router, db *sql.DB) {
	h := &Handler{DB: db}

	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
	router.HandleFunc("/reset-password", h.ResetPassword).Methods("POST")
}
