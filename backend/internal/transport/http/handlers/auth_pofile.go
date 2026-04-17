package handler

import (
	"encoding/json"
	"ethno/internal/transport/http/middleware"
	"net/http"
)

// internal/transport/http/handlers/auth.go
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
    userID := middleware.GetUserID(r)
    if userID == "" {
        http.Error(w, `{"error":"Не авторизован"}`, http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":       userID,
        "username": middleware.GetUsername(r),
        "email":    middleware.GetUserEmail(r),
        "role":     middleware.GetUserRole(r),
    })
}