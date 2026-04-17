package handler

import (
	"encoding/json"
	"ethno/internal/models"
	"net/http"
)

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    user, token, err := h.authService.Login(r.Context(), &req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    token,
        Path:     "/",
        Domain:   h.config.Cookie.Domain,
        MaxAge:   int(h.config.JWT.ExpiryHours * 3600),
        HttpOnly: true,
        Secure:   h.config.Cookie.Secure,
        SameSite: http.SameSiteStrictMode,
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user": user,
        "message": "login successful",
    })
}
