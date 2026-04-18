package models

import (
	"encoding/json"
	"time"
)

type Quest struct {
	ID          string    `json:"id" db:"id"`
	Slug        string    `json:"slug" db:"slug"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CoverURL    string    `json:"cover" db:"cover_url"`
	IsActive    bool      `json:"-" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Steps       []Step    `json:"steps" db:"-"`
}

type Step struct {
	ID          string          `json:"id" db:"id"`
	QuestID     string          `json:"-" db:"quest_id"`
	Order       int             `json:"-" db:"step_order"`
	Type        StepType        `json:"type" db:"step_type"`
	Title       string          `json:"title" db:"title"`
	
	RawContent    json.RawMessage `json:"-" db:"content"`
	OnSuccessRaw  json.RawMessage `json:"-" db:"on_success"`
	
	Content   StepContent   `json:"content" db:"-"`
	OnSuccess *OnSuccess    `json:"onSuccess,omitempty" db:"-"`
}
type StepType string

const (
	StepTypeIntro  StepType = "intro"
	StepTypeQuiz   StepType = "quiz"
	StepTypeBuilder StepType = "builder"
)

type StepContent interface {
	isStepContent()
}

type IntroContent struct {
	Text  string `json:"text"`
	Image string `json:"image,omitempty"`
}
func (IntroContent) isStepContent() {}

type QuizContent struct {
	Question string          `json:"question"`
	Options  []QuizOption    `json:"options"`
}
func (QuizContent) isStepContent() {}

type QuizOption struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Correct bool   `json:"correct"`
}

type BuilderContent struct {
    Base      string           `json:"base"`
    Goal      string           `json:"goal"`
    Required  []string         `json:"required,omitempty"`
    Positions []string         `json:"positions,omitempty"`
    Patterns  []PatternItem    `json:"patterns"`
}

type PatternItem struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    File   string `json:"file"`
    Symbol string `json:"symbol"`
}

func (BuilderContent) isStepContent() {}

type OnSuccess struct {
	Unlock []string       `json:"unlock,omitempty"`
	Reward *QuestReward   `json:"reward,omitempty"`
}

type QuestReward struct {
	Badge string `json:"badge,omitempty"`
	Points int   `json:"points,omitempty"`
}

type UserProgress struct {
	UserID         string     `json:"user_id" db:"user_id"`
	QuestID        string     `json:"quest_id" db:"quest_id"`
	CurrentStepID  string     `json:"current_step_id" db:"current_step_id"`
	CompletedSteps []string   `json:"completed_steps" db:"completed_steps"`
	Status         string     `json:"status" db:"status"`
	StartedAt      time.Time  `json:"started_at" db:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

type MapMatchContent struct {
	Goal    string      `json:"goal"`
	Toys    []Toy       `json:"toys"`
	Regions []MapRegion `json:"regions"`
}

type Toy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type MapRegion struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Center      []float64 `json:"center"`
	CorrectToys []string  `json:"correct_toys"`
}

func (MapMatchContent) isStepContent() {}
