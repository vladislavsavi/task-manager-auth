package handlers

import (
	"encoding/json"
	"net/http"
)

// ResetPasswordRequest - структура для сброса пароля
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

// ResetPassword - сброс пароля пользователя
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.NewPassword == "" {
		http.Error(w, "Email and new password are required", http.StatusBadRequest)
		return
	}

	result, err := h.DB.Exec("UPDATE users SET password = $1 WHERE email = $2", req.NewPassword, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successfully"})
}
