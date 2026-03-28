package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/heracle/pt.heracle.fit.go/internal/models"
)

type MealRepo struct {
	pool *pgxpool.Pool
}

func NewMealRepo(pool *pgxpool.Pool) *MealRepo {
	return &MealRepo{pool: pool}
}

func (r *MealRepo) Create(ctx context.Context, userID, mealType, date, mealTime string, data json.RawMessage) (*models.Meal, error) {
	var m models.Meal
	err := r.pool.QueryRow(ctx,
		`INSERT INTO meals (id, user_id, meal_type, date, time, data, created_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW())
		 RETURNING id, user_id, meal_type, date, time, data, created_at`,
		userID, mealType, date, mealTime, data).Scan(
		&m.ID, &m.UserID, &m.MealType, &m.Date, &m.Time, &m.Data, &m.CreatedAt,
	)
	return &m, err
}

func (r *MealRepo) FindByUserAndDate(ctx context.Context, userID, date string) ([]models.Meal, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, meal_type, date, time, data, created_at
		 FROM meals WHERE user_id = $1 AND date = $2 ORDER BY time ASC`, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var m models.Meal
		if err := rows.Scan(&m.ID, &m.UserID, &m.MealType, &m.Date, &m.Time, &m.Data, &m.CreatedAt); err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	if meals == nil {
		meals = []models.Meal{}
	}
	return meals, nil
}

func (r *MealRepo) FindByUsersAndDate(ctx context.Context, userIDs []string, date string) ([]models.Meal, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, meal_type, date, time, data, created_at
		 FROM meals WHERE user_id = ANY($1) AND date = $2`, userIDs, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var m models.Meal
		if err := rows.Scan(&m.ID, &m.UserID, &m.MealType, &m.Date, &m.Time, &m.Data, &m.CreatedAt); err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	if meals == nil {
		meals = []models.Meal{}
	}
	return meals, nil
}
