package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"ethno/internal/models"
)

type QuestRepository struct {
	db *sql.DB
}

func NewQuestRepository(db *sql.DB) *QuestRepository {
	return &QuestRepository{db: db}
}

func (r *QuestRepository) CreateQuest(ctx context.Context, q *models.Quest, steps []models.Step) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO quests (slug, title, description, cover_url, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`, q.Slug, q.Title, q.Description, q.CoverURL, q.IsActive).Scan(&q.ID, &q.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert quest: %w", err)
	}

	for _, step := range steps {
		step.QuestID = q.ID
		if err := r.createStepTx(ctx, tx, &step); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *QuestRepository) GetBySlug(ctx context.Context, slug string) (*models.Quest, error) {
	var q models.Quest

	err := r.db.QueryRowContext(ctx, `
		SELECT id, slug, title, description, cover_url, is_active, created_at
		FROM quests 
		WHERE slug = $1 AND is_active = true
	`, slug).Scan(&q.ID, &q.Slug, &q.Title, &q.Description, &q.CoverURL, &q.IsActive, &q.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("query quest: %w", err)
	}

	steps, err := r.getStepsByQuestID(ctx, q.ID)
	if err != nil {
		return nil, fmt.Errorf("load steps: %w", err)
	}
	q.Steps = steps

	return &q, nil
}

func (r *QuestRepository) ListActive(ctx context.Context) ([]models.Quest, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, slug, title, description, cover_url, created_at
		FROM quests 
		WHERE is_active = true 
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("query list: %w", err)
	}
	defer rows.Close()

	var quests []models.Quest
	for rows.Next() {
		var q models.Quest
		err := rows.Scan(&q.ID, &q.Slug, &q.Title, &q.Description, &q.CoverURL, &q.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		quests = append(quests, q)
	}
	return quests, rows.Err()
}

func (r *QuestRepository) UpdateQuest(ctx context.Context, q *models.Quest) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE quests 
		SET title = $1, description = $2, cover_url = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5
	`, q.Title, q.Description, q.CoverURL, q.IsActive, q.ID)

	if err != nil {
		return fmt.Errorf("update quest: %w", err)
	}
	return nil
}

func (r *QuestRepository) DeleteQuest(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE quests SET is_active = false, updated_at = NOW() WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("soft delete quest: %w", err)
	}
	return nil
}

func (r *QuestRepository) createStepTx(ctx context.Context, tx *sql.Tx, s *models.Step) error {
	content, err := json.Marshal(s.Content)
	if err != nil {
		return fmt.Errorf("marshal content: %w", err)
	}

	var onSuccessJSON []byte
	if s.OnSuccess != nil {
		onSuccessJSON, err = json.Marshal(s.OnSuccess)
		if err != nil {
			return fmt.Errorf("marshal on_success: %w", err)
		}
	}

	return tx.QueryRowContext(ctx, `
		INSERT INTO quest_steps 
			(quest_id, step_id, step_order, step_type, title, content, on_success)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, s.QuestID, s.ID, s.Order, s.Type, s.Title, content, onSuccessJSON).Scan(&s.ID)
}

func (r *QuestRepository) getStepsByQuestID(ctx context.Context, questID string) ([]models.Step, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, step_id, step_order, step_type, title, content, on_success
		FROM quest_steps 
		WHERE quest_id = $1 
		ORDER BY step_order ASC
	`, questID)
	if err != nil {
		return nil, fmt.Errorf("query steps: %w", err)
	}
	defer rows.Close()

	var steps []models.Step
	for rows.Next() {
		var s models.Step
		var contentJSON, onSuccessJSON []byte

		err := rows.Scan(&s.ID, &s.QuestID, &s.Order, &s.Type, &s.Title, &contentJSON, &onSuccessJSON)
		if err != nil {
			return nil, fmt.Errorf("scan step: %w", err)
		}

		s.RawContent = contentJSON
		if len(onSuccessJSON) > 0 {
			s.OnSuccessRaw = onSuccessJSON
		}
		steps = append(steps, s)
	}
	return steps, rows.Err()
}

func (r *QuestRepository) UpdateStep(ctx context.Context, questID, stepID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0, len(updates)+2)
	argIdx := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIdx))
		args = append(args, value)
		argIdx++
	}
	args = append(args, questID, stepID)

	query := fmt.Sprintf(`
		UPDATE quest_steps 
		SET %s, updated_at = NOW() 
		WHERE quest_id = $%d AND step_id = $%d
	`, 
		fmt.Sprintf("%s, updated_at = NOW()", setClauses), 
		argIdx, argIdx+1)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update step: %w", err)
	}
	return nil
}

func (r *QuestRepository) DeleteStep(ctx context.Context, questID, stepID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM quest_steps WHERE quest_id = $1 AND step_id = $2
	`, questID, stepID)
	if err != nil {
		return fmt.Errorf("delete step: %w", err)
	}
	return nil
}

func (r *QuestRepository) GetProgress(ctx context.Context, userID, questID string) (*models.UserProgress, error) {
	return nil, sql.ErrNoRows // заглушка
}

func (r *QuestRepository) UpsertProgress(ctx context.Context, p *models.UserProgress) error {
	return nil
}

func (r *QuestRepository) GrantReward(ctx context.Context, userID, questID, rewardType, rewardKey string, metadata json.RawMessage) error {
	return nil
}