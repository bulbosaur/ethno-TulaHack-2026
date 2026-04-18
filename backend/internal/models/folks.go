package models

import (
	"encoding/json"
	"time"
)

type Folk struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Lat         *float64        `json:"lat,omitempty"`
	Lon         *float64        `json:"lon,omitempty"`
	Title       *string         `json:"title,omitempty"`
	
	RawSummary  json.RawMessage `json:"-" db:"summary"`
	Summary     *FolkContent    `json:"summary,omitempty" db:"-"`
	
	CreatedAt   time.Time       `json:"created_at"`
}

type FolkContent struct {
	Text        string   `json:"text"`
	Images      []string `json:"images,omitempty"`
	VideoURL    string   `json:"video_url,omitempty"`
	ExternalURL string   `json:"external_url,omitempty"`
}

type CreateFolkDTO struct {
	Name     string
	Lat      *float64
	Lon      *float64
	Title    *string
	Summary  json.RawMessage
}

type UpdateFolkDTO struct {
	Name     string
	Lat      *float64
	Lon      *float64
	Title    *string
	Summary  json.RawMessage
}