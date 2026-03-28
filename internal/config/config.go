package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string

	AdminUsername string
	AdminPassword string

	// Firebase
	FirebaseProjectID   string
	FirebaseClientEmail string
	FirebasePrivateKey  string

	// AI Keys
	OpenAIKey     string
	GeminiKey     string
	HuggingFaceKey string

	// Per-feature AI config
	FoodAnalyseProvider string
	FoodAnalyseModel    string
	DietSuggestionProvider string
	DietSuggestionModel    string
	SleepInsightProvider string
	SleepInsightModel    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		Port:        getEnv("PORT", "3000"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "change-me"),

		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin_password_123"),

		FirebaseProjectID:   getEnv("FIREBASE_PROJECT_ID", ""),
		FirebaseClientEmail: getEnv("FIREBASE_CLIENT_EMAIL", ""),
		FirebasePrivateKey:  getEnv("FIREBASE_PRIVATE_KEY", ""),

		OpenAIKey:     getEnv("OPENAI_API", ""),
		GeminiKey:     getEnv("GEMINI_API", ""),
		HuggingFaceKey: getEnv("HUGGINGFACE_API", ""),

		FoodAnalyseProvider:    getEnv("FOOD_ANALYSE_PROVIDER", "openai"),
		FoodAnalyseModel:       getEnv("FOOD_ANALYSE_MODEL", ""),
		DietSuggestionProvider: getEnv("DIET_SUGGESTION_PROVIDER", "openai"),
		DietSuggestionModel:    getEnv("DIET_SUGGESTION_MODEL", ""),
		SleepInsightProvider:   getEnv("SLEEP_INSIGHT_PROVIDER", "openai"),
		SleepInsightModel:      getEnv("SLEEP_INSIGHT_MODEL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
