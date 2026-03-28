package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutLogRepo struct {
	pool *pgxpool.Pool
}

func NewWorkoutLogRepo(pool *pgxpool.Pool) *WorkoutLogRepo {
	return &WorkoutLogRepo{pool: pool}
}

func (r *WorkoutLogRepo) Create(ctx context.Context, userID string, sessionID *int, logData json.RawMessage, notes *string) (map[string]interface{}, error) {
	var id int
	var uid string
	var sid *int
	var ld json.RawMessage
	var n, ptReview *string
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`INSERT INTO workout_logs (user_id, session_id, log_data, notes, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())
		 RETURNING id, user_id, session_id, log_data, notes, pt_review, created_at, updated_at`,
		userID, sessionID, logData, notes).Scan(&id, &uid, &sid, &ld, &n, &ptReview, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "userId": uid, "sessionId": sid, "logData": ld,
		"notes": n, "ptReview": ptReview, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *WorkoutLogRepo) FindByIDAndUser(ctx context.Context, id int, userID string) (map[string]interface{}, error) {
	var lid int
	var uid string
	var sid *int
	var ld json.RawMessage
	var notes, ptReview *string
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, session_id, log_data, notes, pt_review, created_at, updated_at
		 FROM workout_logs WHERE id = $1 AND user_id = $2`, id, userID).Scan(
		&lid, &uid, &sid, &ld, &notes, &ptReview, &createdAt, &updatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": lid, "userId": uid, "sessionId": sid, "logData": ld,
		"notes": notes, "ptReview": ptReview, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *WorkoutLogRepo) FindByID(ctx context.Context, id int) (map[string]interface{}, error) {
	var lid int
	var uid string
	var sid *int
	var ld json.RawMessage
	var notes, ptReview *string
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, session_id, log_data, notes, pt_review, created_at, updated_at
		 FROM workout_logs WHERE id = $1`, id).Scan(
		&lid, &uid, &sid, &ld, &notes, &ptReview, &createdAt, &updatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": lid, "userId": uid, "sessionId": sid, "logData": ld,
		"notes": notes, "ptReview": ptReview, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *WorkoutLogRepo) FindByUser(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, session_id, log_data, notes, pt_review, created_at, updated_at
		 FROM workout_logs WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var uid string
		var sid *int
		var ld json.RawMessage
		var notes, ptReview *string
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&id, &uid, &sid, &ld, &notes, &ptReview, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id": id, "userId": uid, "sessionId": sid, "logData": ld,
			"notes": notes, "ptReview": ptReview, "createdAt": createdAt, "updatedAt": updatedAt,
		})
	}
	if results == nil {
		results = []map[string]interface{}{}
	}
	return results, nil
}

func (r *WorkoutLogRepo) Update(ctx context.Context, id int, data map[string]interface{}) (map[string]interface{}, error) {
	var lid int
	var uid string
	var sid *int
	var ld json.RawMessage
	var notes, ptReview *string
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`UPDATE workout_logs SET
		   session_id = COALESCE($2, session_id),
		   log_data = COALESCE($3, log_data),
		   notes = COALESCE($4, notes),
		   updated_at = NOW()
		 WHERE id = $1
		 RETURNING id, user_id, session_id, log_data, notes, pt_review, created_at, updated_at`,
		id, data["sessionId"], data["logData"], data["notes"]).Scan(
		&lid, &uid, &sid, &ld, &notes, &ptReview, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": lid, "userId": uid, "sessionId": sid, "logData": ld,
		"notes": notes, "ptReview": ptReview, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *WorkoutLogRepo) UpdatePTReview(ctx context.Context, id int, review string) (map[string]interface{}, error) {
	var lid int
	var uid string
	var sid *int
	var ld json.RawMessage
	var notes, ptReview *string
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`UPDATE workout_logs SET pt_review = $2, updated_at = NOW() WHERE id = $1
		 RETURNING id, user_id, session_id, log_data, notes, pt_review, created_at, updated_at`,
		id, review).Scan(
		&lid, &uid, &sid, &ld, &notes, &ptReview, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": lid, "userId": uid, "sessionId": sid, "logData": ld,
		"notes": notes, "ptReview": ptReview, "createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *WorkoutLogRepo) Delete(ctx context.Context, id int) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM workout_logs WHERE id = $1`, id)
	return err
}
