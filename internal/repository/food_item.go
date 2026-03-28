package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/heracle/pt.heracle.fit.go/internal/models"
)

type FoodItemRepo struct {
	pool *pgxpool.Pool
}

func NewFoodItemRepo(pool *pgxpool.Pool) *FoodItemRepo {
	return &FoodItemRepo{pool: pool}
}

func (r *FoodItemRepo) Search(ctx context.Context, query string) ([]models.FoodItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, calories, protein, carbs, fat, fiber
		 FROM food_items WHERE name ILIKE '%' || $1 || '%' LIMIT 20`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.FoodItem
	for rows.Next() {
		var f models.FoodItem
		if err := rows.Scan(&f.ID, &f.Name, &f.Calories, &f.Protein, &f.Carbs, &f.Fat, &f.Fiber); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	if items == nil {
		items = []models.FoodItem{}
	}
	return items, nil
}
