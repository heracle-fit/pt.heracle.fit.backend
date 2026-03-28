package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DietSuggestionRepo struct {
	pool *pgxpool.Pool
}

func NewDietSuggestionRepo(pool *pgxpool.Pool) *DietSuggestionRepo {
	return &DietSuggestionRepo{pool: pool}
}

func (r *DietSuggestionRepo) FindByUserAndDate(ctx context.Context, userID, date string) (map[string]interface{}, error) {
	var id, uid, suggestion, d string
	var suggestedMeal json.RawMessage
	var createdAt time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, suggestion, suggested_meal, date, created_at
		 FROM diet_suggestions WHERE user_id = $1 AND date = $2
		 ORDER BY created_at DESC LIMIT 1`, userID, date).Scan(
		&id, &uid, &suggestion, &suggestedMeal, &d, &createdAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "suggestion": suggestion, "suggestedMeal": suggestedMeal,
		"date": d, "createdAt": createdAt,
	}, nil
}

func (r *DietSuggestionRepo) Upsert(ctx context.Context, userID, date, suggestion string, suggestedMeal json.RawMessage) (map[string]interface{}, error) {
	var id, uid, s, d string
	var sm json.RawMessage
	var createdAt time.Time

	err := r.pool.QueryRow(ctx,
		`INSERT INTO diet_suggestions (id, user_id, date, suggestion, suggested_meal, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW(), NOW())
		 ON CONFLICT (user_id, date) DO UPDATE SET
		   suggestion = $3,
		   suggested_meal = $4,
		   updated_at = NOW()
		 RETURNING id, user_id, suggestion, suggested_meal, date, created_at`,
		userID, date, suggestion, suggestedMeal).Scan(&id, &uid, &s, &sm, &d, &createdAt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "suggestion": s, "suggestedMeal": sm,
		"date": d, "createdAt": createdAt,
	}, nil
}
