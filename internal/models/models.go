package models

import (
	"encoding/json"
	"time"
)

// ── User ────────────────────────────────────────────────────────────────────────

type User struct {
	ID                string     `json:"id"`
	Username          string     `json:"username"`
	Name              string     `json:"name"`
	Bio               *string    `json:"bio"`
	AvatarURL         *string    `json:"avatarUrl"`
	Email             string     `json:"email"`
	GoogleAccessToken *string    `json:"googleAccessToken,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
	Role              string     `json:"role,omitempty"` // derived, not in DB
}

type UserProfile struct {
	ID                      string     `json:"id"`
	UserID                  string     `json:"userId"`
	Age                     *int       `json:"age"`
	Gender                  *string    `json:"gender"`
	HeightCm                *float64   `json:"heightCm"`
	HeightFt                *float64   `json:"heightFt"`
	WeightKg                *float64   `json:"weightKg"`
	WeightLbs               *float64   `json:"weightLbs"`
	BodyType                *string    `json:"bodyType"`
	Goal                    *string    `json:"goal"`
	FitnessLevel            *string    `json:"fitnessLevel"`
	BMI                     *float64   `json:"bmi"`
	MaintenanceCalories     *int       `json:"maintenanceCalories"`
	GoalWeightKg            *float64   `json:"goalWeightKg"`
	GoalWeightLbs           *float64   `json:"goalWeightLbs"`
	TargetCalories          *int       `json:"targetCalories"`
	TargetProtein           *float64   `json:"targetProtein"`
	TargetCarbs             *float64   `json:"targetCarbs"`
	TargetFat               *float64   `json:"targetFat"`
	TargetFiber             *float64   `json:"targetFiber"`
	WorkoutFrequencyPerWeek *int       `json:"workoutFrequencyPerWeek"`
	PreferredWorkoutType    *string    `json:"preferredWorkoutType"`
	AvailableDays           []string   `json:"availableDays"`
	PreferredWorkoutTime    *string    `json:"preferredWorkoutTime"`
	SessionDurationMins     *int       `json:"sessionDurationMins"`
	Injuries                *string    `json:"injuries"`
	DietaryPreference       *string    `json:"dietaryPreference"`
	DailyWaterLitres        *float64   `json:"dailyWaterLitres"`
	MealsPerDay             *int       `json:"mealsPerDay"`
	CreatedAt               time.Time  `json:"createdAt"`
	UpdatedAt               *time.Time `json:"updatedAt"`
}

// ── Meal ────────────────────────────────────────────────────────────────────────

type Meal struct {
	ID        string          `json:"id"`
	UserID    string          `json:"userId"`
	MealType  string          `json:"mealType"`
	Date      string          `json:"date"`
	Time      string          `json:"time"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"createdAt"`
}

// ── Diet Suggestion ─────────────────────────────────────────────────────────────

type DietSuggestion struct {
	ID            string          `json:"id"`
	UserID        string          `json:"userId"`
	Suggestion    string          `json:"suggestion"`
	SuggestedMeal json.RawMessage `json:"suggestedMeal"`
	Date          string          `json:"date"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

// ── Food Item ───────────────────────────────────────────────────────────────────

type FoodItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

// ── Exercise ────────────────────────────────────────────────────────────────────

type Exercise struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	SecondaryInfo string `json:"secondaryInfo"`
	ExerciseType  string `json:"exerciseType"`
	Image         string `json:"image,omitempty"` // derived
}

// ── Session ─────────────────────────────────────────────────────────────────────

type Session struct {
	ID                   int             `json:"id"`
	UserID               string          `json:"-"`
	Name                 string          `json:"name"`
	Category             json.RawMessage `json:"category"`
	SessionData          json.RawMessage `json:"sessionData"`
	ExerciseImageBaseURL string          `json:"exerciseImageBaseUrl,omitempty"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
}

// ── Workout Log ─────────────────────────────────────────────────────────────────

type WorkoutLog struct {
	ID                   int             `json:"id"`
	UserID               string          `json:"userId"`
	SessionID            *int            `json:"sessionId"`
	LogData              json.RawMessage `json:"logData"`
	Notes                *string         `json:"notes"`
	PTReview             *string         `json:"ptReview"`
	ExerciseImageBaseURL string          `json:"exerciseImageBaseUrl,omitempty"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
}

// ── Sleep Cycle ─────────────────────────────────────────────────────────────────

type SleepCycle struct {
	ID          string          `json:"id"`
	UserID      string          `json:"userId"`
	SleepData   json.RawMessage `json:"sleepData"`
	Insight     *string         `json:"insight"`
	InsightDate *string         `json:"insightDate"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// ── Trainer ─────────────────────────────────────────────────────────────────────

type Trainer struct {
	ID             string     `json:"id"`
	UserID         string     `json:"userId"`
	Specialization *string    `json:"specialization"`
	Experience     *int       `json:"experience"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type TrainerClient struct {
	ID         string    `json:"id"`
	TrainerID  string    `json:"trainerId"`
	ClientID   string    `json:"clientId"`
	AssignedAt time.Time `json:"assignedAt"`
}

// ── Workout Split ───────────────────────────────────────────────────────────────

type WorkoutSplit struct {
	ID        string          `json:"id"`
	TrainerID string          `json:"trainerId"`
	UserID    string          `json:"userId"`
	SplitData json.RawMessage `json:"splitData"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// ── JWT Claims ──────────────────────────────────────────────────────────────────

type JWTClaims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Username string `json:"username,omitempty"`
}
