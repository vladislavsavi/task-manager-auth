package handlers

import (
	"encoding/json"
	"net/http"
)

// RegisterRequest - структура для регистрации
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register - регистрация нового пользователя
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var userID int
	err := h.DB.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", req.Email, req.Password).Scan(&userID)
	if err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": userID, "email": req.Email})
}
