package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// LoginRequest - структура для входа
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login - вход пользователя
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var userID int
	var storedPassword string
	err := h.DB.QueryRow("SELECT id, password FROM users WHERE email = $1", req.Email).Scan(&userID, &storedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Простая проверка пароля (в реальном приложении использовать bcrypt)
	if req.Password != storedPassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": userID, "email": req.Email})
}
