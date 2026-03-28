package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/heracle/pt.heracle.fit.go/internal/models"
)

type ExerciseRepo struct {
	pool *pgxpool.Pool
}

func NewExerciseRepo(pool *pgxpool.Pool) *ExerciseRepo {
	return &ExerciseRepo{pool: pool}
}

func (r *ExerciseRepo) FindAll(ctx context.Context) ([]models.Exercise, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, secondary_info, exercise_type FROM exercises ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var e models.Exercise
		if err := rows.Scan(&e.ID, &e.Name, &e.SecondaryInfo, &e.ExerciseType); err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	if exercises == nil {
		exercises = []models.Exercise{}
	}
	return exercises, nil
}
