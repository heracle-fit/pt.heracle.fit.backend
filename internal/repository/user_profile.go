package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserProfileRepo struct {
	pool *pgxpool.Pool
}

func NewUserProfileRepo(pool *pgxpool.Pool) *UserProfileRepo {
	return &UserProfileRepo{pool: pool}
}

func (r *UserProfileRepo) FindByUserID(ctx context.Context, userID string) (map[string]interface{}, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, age, gender, height_cm, height_ft, weight_kg, weight_lbs, body_type,
		        goal, fitness_level, bmi, maintenance_calories,
		        goal_weight_kg, goal_weight_lbs, target_calories, target_protein, target_carbs, target_fat, target_fiber,
		        workout_frequency_per_week, preferred_workout_type, available_days, preferred_workout_time, session_duration_mins,
		        injuries, dietary_preference, daily_water_litres, meals_per_day,
		        created_at, updated_at
		 FROM user_profiles WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var id, uid string
	var age, maintenanceCalories, targetCalories, workoutFreq, sessionDur, mealsPerDay *int
	var gender, bodyType, goal, fitnessLevel, prefWorkoutType, prefWorkoutTime, injuries, dietPref *string
	var heightCm, heightFt, weightKg, weightLbs, bmi, goalWeightKg, goalWeightLbs *float64
	var targetProtein, targetCarbs, targetFat, targetFiber, dailyWater *float64
	var availableDays []string
	var createdAt time.Time
	var updatedAt *time.Time

	err = rows.Scan(
		&id, &uid, &age, &gender, &heightCm, &heightFt, &weightKg, &weightLbs, &bodyType,
		&goal, &fitnessLevel, &bmi, &maintenanceCalories,
		&goalWeightKg, &goalWeightLbs, &targetCalories, &targetProtein, &targetCarbs, &targetFat, &targetFiber,
		&workoutFreq, &prefWorkoutType, &availableDays, &prefWorkoutTime, &sessionDur,
		&injuries, &dietPref, &dailyWater, &mealsPerDay,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id": id, "userId": uid,
		"age": age, "gender": gender, "heightCm": heightCm, "heightFt": heightFt,
		"weightKg": weightKg, "weightLbs": weightLbs, "bodyType": bodyType,
		"goal": goal, "fitnessLevel": fitnessLevel, "bmi": bmi, "maintenanceCalories": maintenanceCalories,
		"goalWeightKg": goalWeightKg, "goalWeightLbs": goalWeightLbs,
		"targetCalories": targetCalories, "targetProtein": targetProtein, "targetCarbs": targetCarbs,
		"targetFat": targetFat, "targetFiber": targetFiber,
		"workoutFrequencyPerWeek": workoutFreq, "preferredWorkoutType": prefWorkoutType,
		"availableDays": availableDays, "preferredWorkoutTime": prefWorkoutTime,
		"sessionDurationMins": sessionDur, "injuries": injuries,
		"dietaryPreference": dietPref, "dailyWaterLitres": dailyWater, "mealsPerDay": mealsPerDay,
		"createdAt": createdAt, "updatedAt": updatedAt,
	}, nil
}

func (r *UserProfileRepo) GetBodyMetrics(ctx context.Context, userID string) (map[string]interface{}, error) {
	var age, maintenanceCalories, targetCalories *int
	var gender, bodyType, goal *string
	var heightCm, heightFt, weightKg, weightLbs, bmi *float64
	var goalWeightKg, goalWeightLbs, targetProtein, targetCarbs, targetFat, targetFiber *float64
	var updatedAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT age, gender, height_cm, height_ft, weight_kg, weight_lbs, body_type, goal, bmi,
		        maintenance_calories, goal_weight_kg, goal_weight_lbs,
		        target_calories, target_protein, target_carbs, target_fat, target_fiber, updated_at
		 FROM user_profiles WHERE user_id = $1`, userID).Scan(
		&age, &gender, &heightCm, &heightFt, &weightKg, &weightLbs, &bodyType, &goal, &bmi,
		&maintenanceCalories, &goalWeightKg, &goalWeightLbs,
		&targetCalories, &targetProtein, &targetCarbs, &targetFat, &targetFiber, &updatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"age": age, "gender": gender, "heightCm": heightCm, "heightFt": heightFt,
		"weightKg": weightKg, "weightLbs": weightLbs, "bodyType": bodyType, "goal": goal, "bmi": bmi,
		"maintenanceCalories": maintenanceCalories,
		"goalWeightKg": goalWeightKg, "goalWeightLbs": goalWeightLbs,
		"targetCalories": targetCalories, "targetProtein": targetProtein, "targetCarbs": targetCarbs,
		"targetFat": targetFat, "targetFiber": targetFiber, "updatedAt": updatedAt,
	}, nil
}

func (r *UserProfileRepo) GetTargets(ctx context.Context, userID string) (map[string]interface{}, error) {
	var targetCalories *int
	var targetProtein, targetCarbs, targetFat, targetFiber *float64

	err := r.pool.QueryRow(ctx,
		`SELECT target_calories, target_protein, target_carbs, target_fat, target_fiber
		 FROM user_profiles WHERE user_id = $1`, userID).Scan(
		&targetCalories, &targetProtein, &targetCarbs, &targetFat, &targetFiber,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"targetCalories": targetCalories, "targetProtein": targetProtein,
		"targetCarbs": targetCarbs, "targetFat": targetFat, "targetFiber": targetFiber,
	}, nil
}

func (r *UserProfileRepo) UpsertBodyMetrics(ctx context.Context, userID string, data map[string]interface{}) (map[string]interface{}, error) {
	// Build a complete upsert with the provided fields
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_profiles (id, user_id, age, gender, height_cm, height_ft, weight_kg, weight_lbs,
		 body_type, goal, goal_weight_kg, goal_weight_lbs, bmi, maintenance_calories,
		 target_calories, target_protein, target_carbs, target_fat, target_fiber, created_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET
		   age = COALESCE($2, user_profiles.age),
		   gender = COALESCE($3, user_profiles.gender),
		   height_cm = COALESCE($4, user_profiles.height_cm),
		   height_ft = COALESCE($5, user_profiles.height_ft),
		   weight_kg = COALESCE($6, user_profiles.weight_kg),
		   weight_lbs = COALESCE($7, user_profiles.weight_lbs),
		   body_type = COALESCE($8, user_profiles.body_type),
		   goal = COALESCE($9, user_profiles.goal),
		   goal_weight_kg = COALESCE($10, user_profiles.goal_weight_kg),
		   goal_weight_lbs = COALESCE($11, user_profiles.goal_weight_lbs),
		   bmi = COALESCE($12, user_profiles.bmi),
		   maintenance_calories = COALESCE($13, user_profiles.maintenance_calories),
		   target_calories = COALESCE($14, user_profiles.target_calories),
		   target_protein = COALESCE($15, user_profiles.target_protein),
		   target_carbs = COALESCE($16, user_profiles.target_carbs),
		   target_fat = COALESCE($17, user_profiles.target_fat),
		   target_fiber = COALESCE($18, user_profiles.target_fiber),
		   updated_at = NOW()`,
		userID,
		data["age"], data["gender"], data["heightCm"], data["heightFt"],
		data["weightKg"], data["weightLbs"], data["bodyType"], data["goal"],
		data["goalWeightKg"], data["goalWeightLbs"],
		data["bmi"], data["maintenanceCalories"],
		data["targetCalories"], data["targetProtein"], data["targetCarbs"], data["targetFat"], data["targetFiber"],
	)
	if err != nil {
		return nil, err
	}
	return r.GetBodyMetrics(ctx, userID)
}

func (r *UserProfileRepo) UpsertTargets(ctx context.Context, userID string, data map[string]interface{}) (map[string]interface{}, error) {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_profiles (id, user_id, target_calories, target_protein, target_carbs, target_fat, target_fiber, created_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET
		   target_calories = COALESCE($2, user_profiles.target_calories),
		   target_protein = COALESCE($3, user_profiles.target_protein),
		   target_carbs = COALESCE($4, user_profiles.target_carbs),
		   target_fat = COALESCE($5, user_profiles.target_fat),
		   target_fiber = COALESCE($6, user_profiles.target_fiber),
		   updated_at = NOW()`,
		userID, data["targetCalories"], data["targetProtein"], data["targetCarbs"], data["targetFat"], data["targetFiber"],
	)
	if err != nil {
		return nil, err
	}

	var targetCalories *int
	var targetProtein, targetCarbs, targetFat, targetFiber *float64
	var updatedAt *time.Time
	err = r.pool.QueryRow(ctx,
		`SELECT target_calories, target_protein, target_carbs, target_fat, target_fiber, updated_at
		 FROM user_profiles WHERE user_id = $1`, userID).Scan(
		&targetCalories, &targetProtein, &targetCarbs, &targetFat, &targetFiber, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"targetCalories": targetCalories, "targetProtein": targetProtein,
		"targetCarbs": targetCarbs, "targetFat": targetFat, "targetFiber": targetFiber,
		"updatedAt": updatedAt,
	}, nil
}

func (r *UserProfileRepo) GetOnboardingStatus(ctx context.Context, userID string) (map[string]interface{}, error) {
	profile, err := r.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	bodyMetricsNeeded := true
	dietDataNeeded := true

	if profile != nil {
		bodyMetricsNeeded = profile["age"] == nil || profile["gender"] == nil ||
			profile["heightCm"] == nil || profile["weightKg"] == nil || profile["goal"] == nil
		dietDataNeeded = profile["dietaryPreference"] == nil || profile["mealsPerDay"] == nil
	}

	return map[string]interface{}{
		"bodyMetricsNeeded": bodyMetricsNeeded,
		"dietDataNeeded":    dietDataNeeded,
	}, nil
}

func (r *UserProfileRepo) GetDietPreferences(ctx context.Context, userID string) (map[string]interface{}, error) {
	var dietPref *string
	var dailyWater *float64
	var mealsPerDay *int
	var updatedAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT dietary_preference, daily_water_litres, meals_per_day, updated_at
		 FROM user_profiles WHERE user_id = $1`, userID).Scan(
		&dietPref, &dailyWater, &mealsPerDay, &updatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"dietaryPreference": dietPref, "dailyWaterLitres": dailyWater, "mealsPerDay": mealsPerDay, "updatedAt": updatedAt,
	}, nil
}

func (r *UserProfileRepo) UpsertDietPreferences(ctx context.Context, userID string, data map[string]interface{}) (map[string]interface{}, error) {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_profiles (id, user_id, dietary_preference, daily_water_litres, meals_per_day, created_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET
		   dietary_preference = COALESCE($2, user_profiles.dietary_preference),
		   daily_water_litres = COALESCE($3, user_profiles.daily_water_litres),
		   meals_per_day = COALESCE($4, user_profiles.meals_per_day),
		   updated_at = NOW()`,
		userID, data["dietaryPreference"], data["dailyWaterLitres"], data["mealsPerDay"],
	)
	if err != nil {
		return nil, err
	}
	return r.GetDietPreferences(ctx, userID)
}

func (r *UserProfileRepo) GetWorkoutPreferences(ctx context.Context, userID string) (map[string]interface{}, error) {
	var id string
	var fitnessLevel, prefWorkoutType, injuries, prefWorkoutTime *string
	var workoutFreq, sessionDur *int
	var availableDays []string
	var updatedAt *time.Time

	err := r.pool.QueryRow(ctx,
		`SELECT id, fitness_level, workout_frequency_per_week, preferred_workout_type,
		        injuries, available_days, preferred_workout_time, session_duration_mins, updated_at
		 FROM user_profiles WHERE user_id = $1`, userID).Scan(
		&id, &fitnessLevel, &workoutFreq, &prefWorkoutType,
		&injuries, &availableDays, &prefWorkoutTime, &sessionDur, &updatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if availableDays == nil {
		availableDays = []string{}
	}
	return map[string]interface{}{
		"id": id, "fitnessLevel": fitnessLevel, "workoutFrequencyPerWeek": workoutFreq,
		"preferredWorkoutType": prefWorkoutType, "injuries": injuries,
		"availableDays": availableDays, "preferredWorkoutTime": prefWorkoutTime,
		"sessionDurationMins": sessionDur, "updatedAt": updatedAt,
	}, nil
}

func (r *UserProfileRepo) UpsertWorkoutPreferences(ctx context.Context, userID string, data map[string]interface{}) (map[string]interface{}, error) {
	var days []string
	if d, ok := data["availableDays"].([]interface{}); ok {
		for _, v := range d {
			if s, ok := v.(string); ok {
				days = append(days, s)
			}
		}
	} else if d, ok := data["availableDays"].([]string); ok {
		days = d
	}
	if days == nil {
		// Fetch existing to preserve
		existing, _ := r.GetWorkoutPreferences(ctx, userID)
		if existing != nil {
			if ed, ok := existing["availableDays"].([]string); ok {
				days = ed
			}
		}
		if days == nil {
			days = []string{}
		}
	}

	daysJSON, _ := json.Marshal(days)

	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_profiles (id, user_id, fitness_level, workout_frequency_per_week, preferred_workout_type,
		 injuries, available_days, preferred_workout_time, session_duration_mins, created_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6::text[], $7, $8, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET
		   fitness_level = COALESCE($2, user_profiles.fitness_level),
		   workout_frequency_per_week = COALESCE($3, user_profiles.workout_frequency_per_week),
		   preferred_workout_type = COALESCE($4, user_profiles.preferred_workout_type),
		   injuries = COALESCE($5, user_profiles.injuries),
		   available_days = COALESCE($6::text[], user_profiles.available_days),
		   preferred_workout_time = COALESCE($7, user_profiles.preferred_workout_time),
		   session_duration_mins = COALESCE($8, user_profiles.session_duration_mins),
		   updated_at = NOW()`,
		userID, data["fitnessLevel"], data["workoutFrequencyPerWeek"], data["preferredWorkoutType"],
		data["injuries"], string(daysJSON), data["preferredWorkoutTime"], data["sessionDurationMins"],
	)
	if err != nil {
		return nil, err
	}
	return r.GetWorkoutPreferences(ctx, userID)
}
