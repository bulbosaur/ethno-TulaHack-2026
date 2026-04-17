package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"ethno/internal/models"
	"ethno/internal/repository"
)

func GetRandomFolksHandler(repo *repository.FolkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 1
		if l := r.URL.Query().Get("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		folks, err := repo.GetRandom(r.Context(), limit)
		if err != nil {
			log.Printf("GetRandomFolksHandler: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to load regions"}`))
			return
		}

		if folks == nil {
			folks = []models.Folk{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(folks)
	}
}
