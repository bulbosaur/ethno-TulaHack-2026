package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	f "ethno/internal/folks"
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

		for i := range folks {
			if folks[i].RawSummary != nil {
				content, err := f.ParseFolkContent(folks[i].RawSummary)
				if err != nil {
					log.Printf("Failed to parse folk content for %s: %v", folks[i].ID, err)
					continue
				}
				folks[i].Summary = content
				folks[i].RawSummary = nil
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(folks); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}

func GetFolkByIDHandler(repo *repository.FolkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		
		folk, err := repo.GetByID(r.Context(), id)
		if err != nil {
			log.Printf("GetFolkByIDHandler: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to load region"}`))
			return
		}

		if folk == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "region not found"}`))
			return
		}

		if folk.RawSummary != nil {
			content, err := f.ParseFolkContent(folk.RawSummary)
			if err != nil {
				log.Printf("Failed to parse folk content for %s: %v", id, err)
			} else {
				folk.Summary = content
				folk.RawSummary = nil
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(folk); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}

func GetFolksHandler(repo *repository.FolkRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		folks, err := repo.GetAll(r.Context())
		if err != nil {
			log.Printf("GetFolksHandler: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "failed to load regions"}`))
			return
		}

		if folks == nil {
			folks = []models.Folk{}
		}

		for i := range folks {
			if folks[i].RawSummary != nil {
				content, err := f.ParseFolkContent(folks[i].RawSummary)
				if err != nil {
					log.Printf("Failed to parse folk content for %s: %v", folks[i].ID, err)
					continue
				}
				folks[i].Summary = content
				folks[i].RawSummary = nil
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(folks); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}
