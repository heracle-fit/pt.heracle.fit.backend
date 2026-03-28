package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepo struct {
	pool *pgxpool.Pool
}

func NewSessionRepo(pool *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{pool: pool}
}

func (r *SessionRepo) Create(ctx context.Context, userID, name string, category, sessionData json.RawMessage) (map[string]interface{}, error) {
	var id int
	var n string
	var cat, sd json.RawMessage
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`INSERT INTO sessions (user_id, name, category, session_data, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())
		 RETURNING id, name, category, session_data, created_at, updated_at`,
		userID, name, category, sessionData).Scan(&id, &n, &cat, &sd, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": id, "name": n, "category": cat, "sessionData": sd,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *SessionRepo) FindByIDAndUser(ctx context.Context, id int, userID string) (map[string]interface{}, error) {
	var sid int
	var name string
	var cat, sd json.RawMessage
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, name, category, session_data, created_at, updated_at
		 FROM sessions WHERE id = $1 AND user_id = $2`, id, userID).Scan(
		&sid, &name, &cat, &sd, &createdAt, &updatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": sid, "name": name, "category": cat, "sessionData": sd,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *SessionRepo) FindByUser(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, category, session_data, created_at, updated_at
		 FROM sessions WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		var cat, sd json.RawMessage
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&id, &name, &cat, &sd, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id": id, "name": name, "category": cat, "sessionData": sd,
			"createdAt": createdAt, "updatedAt": updatedAt,
		})
	}
	if results == nil {
		results = []map[string]interface{}{}
	}
	return results, nil
}

func (r *SessionRepo) Update(ctx context.Context, id int, data map[string]interface{}) (map[string]interface{}, error) {
	var sid int
	var name string
	var cat, sd json.RawMessage
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`UPDATE sessions SET
		   name = COALESCE($2, name),
		   category = COALESCE($3, category),
		   session_data = COALESCE($4, session_data),
		   updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, name, category, session_data, created_at, updated_at`,
		id, data["name"], data["category"], data["sessionData"]).Scan(
		&sid, &name, &cat, &sd, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": sid, "name": name, "category": cat, "sessionData": sd,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *SessionRepo) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}
