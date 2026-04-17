package handler

import (
	"bytes"
	"encoding/json"
	"ethno/internal/models"

	"io"
	"net/http"
)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	h.logger.Debugf("Register: method=%s, content-type=%s", r.Method, r.Header.Get("Content-Type"))

	bodyBytes, _ := io.ReadAll(r.Body)
    h.logger.Debugf("Register: raw body=%s", string(bodyBytes))
    r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
