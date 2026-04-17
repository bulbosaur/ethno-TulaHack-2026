package handler

import (
	"encoding/json"
	"net/http"

	que "ethno/internal/quest"
	"ethno/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type QuestHandler struct {
	repo   *repository.QuestRepository
	logger *logrus.Logger
}

func NewQuestHandler(repo *repository.QuestRepository, logger *logrus.Logger) *QuestHandler {
	return &QuestHandler{repo: repo, logger: logger}
}

func (h *QuestHandler) RegisterRoutes(r chi.Router) {
	r.Route("/api/quests", func(r chi.Router) {
		r.Get("/", h.listQuests)
		r.Get("/{slug}", h.getQuest)
	})
}

func (h *QuestHandler) listQuests(w http.ResponseWriter, r *http.Request) {
	quests, err := h.repo.ListActive(r.Context())
	if err != nil {
		h.logger.Error("list quests: ", err)
		http.Error(w, `{"error":"server error"}`, http.StatusInternalServerError)
		return
	}

	type preview struct {
		ID    string `json:"id"`
		Slug  string `json:"slug"`
		Title string `json:"title"`
		Cover string `json:"cover"`
	}
	
	out := make([]preview, len(quests))
	for i, q := range quests {
		out[i] = preview{ID: q.ID, Slug: q.Slug, Title: q.Title, Cover: q.CoverURL}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"data": out})
}

func (h *QuestHandler) getQuest(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	
	q, err := h.repo.GetBySlug(r.Context(), slug)
	if err != nil {
		h.logger.WithField("slug", slug).Error("get quest: ", err)
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	for i := range q.Steps {
		if err := que.ParseStepContent(&q.Steps[i]); err != nil {
			h.logger.Warn("parse step: ", err)
			continue
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"data": q})
}
