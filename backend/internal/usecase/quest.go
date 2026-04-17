package usecase

import (
	"context"
	"encoding/json"
	"ethno/internal/models"
)

type QuestUsecase interface {
	ListActive(ctx context.Context) ([]models.Quest, error)
	GetBySlug(ctx context.Context, slug string) (*models.Quest, error)
	ProcessStep(ctx context.Context, req StepProcess) (*StepResult, error)
}

type StepProcess struct {
	QuestSlug string
	UserID    string
	StepID    string
	Action    string
	Payload   json.RawMessage
}

type StepResult struct {
	Success      bool              `json:"success"`
	NextStepID   *string           `json:"next_step_id,omitempty"`
	Unlocked     []string          `json:"unlocked,omitempty"`
	Reward       *models.QuestReward `json:"reward,omitempty"`
	Completed    bool              `json:"quest_completed,omitempty"`
	Message      string            `json:"message,omitempty"`
}