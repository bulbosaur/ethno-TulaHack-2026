package handler

import (
	"encoding/json"
	quests "ethno/internal/quest"
	"ethno/internal/usecase"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type QuestHandler struct {
	uc     usecase.QuestUsecase
	logger *logrus.Logger
}

func NewQuestHandler(uc usecase.QuestUsecase, logger *logrus.Logger) *QuestHandler {
	return &QuestHandler{uc: uc, logger: logger}
}

func (h *QuestHandler) RegisterRoutes(r chi.Router) {
	r.Get("/quests", h.listQuests)
	r.Get("/quests/{slug}", h.getQuest)
	r.Post("/quests/{slug}/progress", h.updateProgress)
}

func (h *QuestHandler) listQuests(w http.ResponseWriter, r *http.Request) {
	quests, err := h.uc.ListActive(r.Context())
	if err != nil {
		h.logger.Error("failed to list quests: ", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	type questPreview struct {
		ID          string `json:"id"`
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
		CoverURL    string `json:"cover"`
	}
	previews := make([]questPreview, len(quests))
	for i, q := range quests {
		previews[i] = questPreview{
			ID: q.ID, Slug: q.Slug, Title: q.Title,
			Description: q.Description, CoverURL: q.CoverURL,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": previews,
	})
}

func (h *QuestHandler) getQuest(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	
	quest, err := h.uc.GetBySlug(r.Context(), slug)
	if err != nil {
		h.logger.WithField("slug", slug).Error("quest not found: ", err)
		http.Error(w, `{"error":"quest not found"}`, http.StatusNotFound)
		return
	}

	for i := range quest.Steps {
		if err := quests.ParseStepContent(&quest.Steps[i]); err != nil {
			h.logger.Error("failed to parse step content: ", err)
			continue
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": quest,
	})
}

type progressRequest struct {
	StepID   string          `json:"step_id"`
	Action   string          `json:"action"`
	Payload  json.RawMessage `json:"payload,omitempty"`
	UserID   string          `json:"user_id"`
}

func (h *QuestHandler) updateProgress(w http.ResponseWriter, r *http.Request) {
	var req progressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	result, err := h.uc.ProcessStep(r.Context(), usecase.StepProcess{
		QuestSlug: chi.URLParam(r, "slug"),
		UserID:    req.UserID,
		StepID:    req.StepID,
		Action:    req.Action,
		Payload:   req.Payload,
	})
	if err != nil {
		h.logger.Error("process step failed: ", err)
		http.Error(w, `{"error":"processing failed"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": result,
	})
}