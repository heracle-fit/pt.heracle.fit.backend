package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SleepRepo struct {
	pool *pgxpool.Pool
}

func NewSleepRepo(pool *pgxpool.Pool) *SleepRepo {
	return &SleepRepo{pool: pool}
}

func (r *SleepRepo) FindByUser(ctx context.Context, userID string) (map[string]interface{}, error) {
	var id, uid string
	var sleepData json.RawMessage
	var insight, insightDate *string
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, sleep_data, insight, insight_date, created_at, updated_at
		 FROM sleep_cycles WHERE user_id = $1 LIMIT 1`, userID).Scan(
		&id, &uid, &sleepData, &insight, &insightDate, &createdAt, &updatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "userId": uid, "sleepData": sleepData,
		"insight": insight, "insightDate": insightDate,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *SleepRepo) Upsert(ctx context.Context, userID string, sleepData json.RawMessage, insight *string, insightDate *string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sleep_cycles (id, user_id, sleep_data, insight, insight_date, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())
		 ON CONFLICT ((user_id)) DO NOTHING`,
		userID, sleepData, insight, insightDate)

	// ON CONFLICT doesn't work well here since user_id isn't unique.
	// Instead, use the find-then-update/create pattern
	return err
}

func (r *SleepRepo) UpdateSleepData(ctx context.Context, id string, sleepData json.RawMessage) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sleep_cycles SET sleep_data = $2, updated_at = NOW() WHERE id = $1`, id, sleepData)
	return err
}

func (r *SleepRepo) UpdateInsight(ctx context.Context, id string, insight string, insightDate string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE sleep_cycles SET insight = $2, insight_date = $3, updated_at = NOW() WHERE id = $1`,
		id, insight, insightDate)
	return err
}

func (r *SleepRepo) Create(ctx context.Context, userID string, sleepData json.RawMessage) (string, error) {
	var id string
	err := r.pool.QueryRow(ctx,
		`INSERT INTO sleep_cycles (id, user_id, sleep_data, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, NOW(), NOW())
		 RETURNING id`, userID, sleepData).Scan(&id)
	return id, err
}

func (r *SleepRepo) CreateWithInsight(ctx context.Context, userID string, sleepData json.RawMessage, insight, insightDate string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO sleep_cycles (id, user_id, sleep_data, insight, insight_date, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())`,
		userID, sleepData, insight, insightDate)
	return err
}
