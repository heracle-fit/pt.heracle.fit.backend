package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SplitRepo struct {
	pool *pgxpool.Pool
}

func NewSplitRepo(pool *pgxpool.Pool) *SplitRepo {
	return &SplitRepo{pool: pool}
}

func (r *SplitRepo) FindByUserID(ctx context.Context, userID string) (map[string]interface{}, error) {
	var id, trainerID, uid string
	var splitData json.RawMessage
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, trainer_id, user_id, split_data, created_at, updated_at
		 FROM workout_splits WHERE user_id = $1`, userID).Scan(
		&id, &trainerID, &uid, &splitData, &createdAt, &updatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "trainerId": trainerID, "userId": uid, "splitData": splitData,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *SplitRepo) Upsert(ctx context.Context, userID, trainerID string, splitData json.RawMessage) (map[string]interface{}, error) {
	var id, tid, uid string
	var sd json.RawMessage
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`INSERT INTO workout_splits (id, trainer_id, user_id, split_data, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		 ON CONFLICT (user_id) DO UPDATE SET
		   trainer_id = $1,
		   split_data = $3,
		   updated_at = NOW()
		 RETURNING id, trainer_id, user_id, split_data, created_at, updated_at`,
		trainerID, userID, splitData).Scan(&id, &tid, &uid, &sd, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "trainerId": tid, "userId": uid, "splitData": sd,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}
