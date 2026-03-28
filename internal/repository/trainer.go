package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/heracle/pt.heracle.fit.go/internal/models"
)

type TrainerRepo struct {
	pool *pgxpool.Pool
}

func NewTrainerRepo(pool *pgxpool.Pool) *TrainerRepo {
	return &TrainerRepo{pool: pool}
}

func (r *TrainerRepo) FindByUserID(ctx context.Context, userID string) (*models.Trainer, error) {
	var t models.Trainer
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, specialization, experience, created_at, updated_at
		 FROM trainers WHERE user_id = $1`, userID).Scan(
		&t.ID, &t.UserID, &t.Specialization, &t.Experience, &t.CreatedAt, &t.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (r *TrainerRepo) FindClientAssignment(ctx context.Context, clientID string) (*models.TrainerClient, error) {
	var tc models.TrainerClient
	err := r.pool.QueryRow(ctx,
		`SELECT id, trainer_id, client_id, assigned_at FROM trainer_clients WHERE client_id = $1`, clientID).Scan(
		&tc.ID, &tc.TrainerID, &tc.ClientID, &tc.AssignedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &tc, err
}

func (r *TrainerRepo) FindClientsByTrainer(ctx context.Context, trainerID string) ([]map[string]interface{}, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT tc.client_id, tc.assigned_at,
		        u.id, u.name, u.email, u.avatar_url,
		        up.goal, up.target_calories
		 FROM trainer_clients tc
		 JOIN users u ON u.id = tc.client_id
		 LEFT JOIN user_profiles up ON up.user_id = tc.client_id
		 WHERE tc.trainer_id = $1`, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var clientID string
		var assignedAt time.Time
		var id, name, email string
		var avatarURL, goal *string
		var targetCalories *int

		if err := rows.Scan(&clientID, &assignedAt, &id, &name, &email, &avatarURL, &goal, &targetCalories); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"id": id, "name": name, "email": email, "avatarUrl": avatarURL,
			"assignedAt": assignedAt, "goal": goal, "targetCalories": targetCalories,
		})
	}
	if results == nil {
		results = []map[string]interface{}{}
	}
	return results, nil
}

func (r *TrainerRepo) AddClient(ctx context.Context, trainerID, clientID string) (*models.TrainerClient, error) {
	var tc models.TrainerClient
	err := r.pool.QueryRow(ctx,
		`INSERT INTO trainer_clients (id, trainer_id, client_id, assigned_at)
		 VALUES (gen_random_uuid(), $1, $2, NOW())
		 RETURNING id, trainer_id, client_id, assigned_at`,
		trainerID, clientID).Scan(&tc.ID, &tc.TrainerID, &tc.ClientID, &tc.AssignedAt)
	return &tc, err
}

func (r *TrainerRepo) RemoveClient(ctx context.Context, clientID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM trainer_clients WHERE client_id = $1`, clientID)
	return err
}

func (r *TrainerRepo) CreateTrainer(ctx context.Context, userID string, specialization *string, experience *int) (map[string]interface{}, error) {
	var id, uid string
	var spec *string
	var exp *int
	var createdAt, updatedAt time.Time

	err := r.pool.QueryRow(ctx,
		`INSERT INTO trainers (id, user_id, specialization, experience, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		 RETURNING id, user_id, specialization, experience, created_at, updated_at`,
		userID, specialization, experience).Scan(&id, &uid, &spec, &exp, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "userId": uid, "specialization": spec, "experience": exp,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *TrainerRepo) GetClientDetails(ctx context.Context, clientID string) (map[string]interface{}, error) {
	var id, name, email string
	var avatarURL *string

	err := r.pool.QueryRow(ctx,
		`SELECT id, name, email, avatar_url FROM users WHERE id = $1`, clientID).Scan(
		&id, &name, &email, &avatarURL)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id": id, "name": name, "email": email, "avatarUrl": avatarURL,
	}, nil
}
