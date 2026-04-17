package models

import "time"

type Folk struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Lat       *float64  `json:"lat,omitempty"`
	Lon       *float64  `json:"lon,omitempty"`
	Title     *string   `json:"title,omitempty"`
	Summary   *string   `json:"summary,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}


type CreateFolkDTO struct {
	Name    string
	Lat     *float64
	Lon     *float64
	Title   *string
	Summary *string
}

type UpdateFolkDTO struct {
	Name    string
	Lat     *float64
	Lon     *float64
	Title   *string
	Summary *string
}
