package repository

import (
	"context"
	"database/sql"
	"ethno/internal/models"
	"fmt"
)

type FolkRepository struct {
	db *sql.DB
}

func NewFolkRepository(db *sql.DB) *FolkRepository {
	return &FolkRepository{db: db}
}

func (r *FolkRepository) Create(ctx context.Context, dto models.CreateFolkDTO) (*models.Folk, error) {
	query := `
        INSERT INTO folks (name, lat, lon, title, summary) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id, created_at
    `

	var folk models.Folk
	err := r.db.QueryRowContext(ctx, query, dto.Name, dto.Lat, dto.Lon, dto.Title, dto.Summary).
		Scan(&folk.ID, &folk.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create folk: %w", err)
	}

	folk.Name = dto.Name
	folk.Lat = dto.Lat
	folk.Lon = dto.Lon
	folk.Title = dto.Title
	folk.RawSummary = dto.Summary

	return &folk, nil
}

func (r *FolkRepository) GetByID(ctx context.Context, id string) (*models.Folk, error) {
	query := `
        SELECT id, name, lat, lon, title, summary, created_at 
        FROM folks 
        WHERE id = $1
    `

	var f models.Folk
	var rawSummary []byte
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&f.ID, &f.Name, &f.Lat, &f.Lon, &f.Title, &rawSummary, &f.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get folk by id: %w", err)
	}

	if rawSummary != nil {
		f.RawSummary = rawSummary
	}

	return &f, nil
}

func (r *FolkRepository) GetRandom(ctx context.Context, limit int) ([]models.Folk, error) {
	query := `
		SELECT id, name, lat, lon, title, summary, created_at 
		FROM folks 
		ORDER BY RANDOM() 
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get random folks: %w", err)
	}
	defer rows.Close()

	var folks []models.Folk
	for rows.Next() {
		var f models.Folk
		var rawSummary []byte
		if err := rows.Scan(&f.ID, &f.Name, &f.Lat, &f.Lon, &f.Title, &rawSummary, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan folk row: %w", err)
		}
		if rawSummary != nil {
			f.RawSummary = rawSummary
		}
		folks = append(folks, f)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return folks, nil
}

func (r *FolkRepository) Update(ctx context.Context, id string, dto models.UpdateFolkDTO) (*models.Folk, error) {
	query := `
        UPDATE folks 
        SET name = $2, lat = $3, lon = $4, title = $5, summary = $6
        WHERE id = $1
        RETURNING id, name, lat, lon, title, summary, created_at
    `

	var f models.Folk
	var rawSummary []byte
	err := r.db.QueryRowContext(ctx, query, id, dto.Name, dto.Lat, dto.Lon, dto.Title, dto.Summary).
		Scan(&f.ID, &f.Name, &f.Lat, &f.Lon, &f.Title, &rawSummary, &f.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update folk: %w", err)
	}

	if rawSummary != nil {
		f.RawSummary = rawSummary
	}

	return &f, nil
}

func (r *FolkRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM folks WHERE id = $1`
	
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete folk: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("folk with id %s not found", id)
	}

	return nil
}

func (r *FolkRepository) List(ctx context.Context) ([]models.Folk, error) {
	query := `
        SELECT id, name, lat, lon, title, summary, created_at 
        FROM folks 
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list folks: %w", err)
	}
	defer rows.Close()

	var folks []models.Folk
	for rows.Next() {
		var f models.Folk
		var rawSummary []byte
		if err := rows.Scan(&f.ID, &f.Name, &f.Lat, &f.Lon, &f.Title, &rawSummary, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan folk row: %w", err)
		}
		if rawSummary != nil {
			f.RawSummary = rawSummary
		}
		folks = append(folks, f)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return folks, nil
}