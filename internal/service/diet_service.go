package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/heracle/pt.heracle.fit.go/internal/ai"
	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

type DietService struct {
	mealRepo       *repository.MealRepo
	profileRepo    *repository.UserProfileRepo
	suggestionRepo *repository.DietSuggestionRepo
	foodItemRepo   *repository.FoodItemRepo
	trainerRepo    *repository.TrainerRepo
	aiRouter       *ai.AIRouter
}

func NewDietService(
	mealRepo *repository.MealRepo,
	profileRepo *repository.UserProfileRepo,
	suggestionRepo *repository.DietSuggestionRepo,
	foodItemRepo *repository.FoodItemRepo,
	trainerRepo *repository.TrainerRepo,
	aiRouter *ai.AIRouter,
) *DietService {
	return &DietService{
		mealRepo: mealRepo, profileRepo: profileRepo,
		suggestionRepo: suggestionRepo, foodItemRepo: foodItemRepo,
		trainerRepo: trainerRepo, aiRouter: aiRouter,
	}
}

func (s *DietService) GetStatus(ctx context.Context, userID, date string) (map[string]interface{}, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.getDailyNutritionalStatus(ctx, userID, date)
}

func (s *DietService) GetTodayDiet(ctx context.Context, userID string) (interface{}, error) {
	today := time.Now().Format("2006-01-02")

	existing, err := s.suggestionRepo.FindByUserAndDate(ctx, userID, today)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	if existing != nil {
		return existing, nil
	}

	generated, err := s.generateDietSuggestion(ctx, userID)
	if err != nil || generated == nil {
		return nil, nil
	}
	return generated, nil
}

func (s *DietService) SearchFood(ctx context.Context, query string) (interface{}, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []interface{}{}, nil
	}
	items, err := s.foodItemRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return items, nil
}

func (s *DietService) GetMealsByDate(ctx context.Context, userID, date string) (interface{}, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	meals, err := s.mealRepo.FindByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return meals, nil
}

func (s *DietService) GetDietPreferences(ctx context.Context, userID string) (map[string]interface{}, error) {
	prefs, err := s.profileRepo.GetDietPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return prefs, nil
}

func (s *DietService) SaveDietPreferences(ctx context.Context, userID string, body map[string]interface{}) (map[string]interface{}, error) {
	result, err := s.profileRepo.UpsertDietPreferences(ctx, userID, body)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return result, nil
}

type LogMealRequest struct {
	MealType string          `json:"mealType"`
	Date     string          `json:"date"`
	Time     string          `json:"time"`
	Data     json.RawMessage `json:"data"`
}

func (s *DietService) LogMeal(ctx context.Context, userID string, req LogMealRequest) (map[string]interface{}, error) {
	meal, err := s.mealRepo.Create(ctx, userID, req.MealType, req.Date, req.Time, req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to log meal")
	}

	// Trigger AI suggestion in background
	go func() {
		suggestion, err := s.generateDietSuggestion(context.Background(), userID)
		if err != nil {
			log.Printf("Background diet suggestion failed: %v", err)
		} else {
			log.Printf("Background diet suggestion generated: %v", suggestion != nil)
		}
	}()

	return map[string]interface{}{
		"id":        meal.ID,
		"userId":    meal.UserID,
		"mealType":  meal.MealType,
		"date":      meal.Date,
		"time":      meal.Time,
		"data":      meal.Data,
		"createdAt": meal.CreatedAt,
	}, nil
}

func (s *DietService) AnalyseFood(ctx context.Context, description *string, imageData []byte, mimeType string) (interface{}, error) {
	if len(imageData) == 0 && (description == nil || *description == "") {
		return nil, ErrBadRequest("at least one of image or description must be provided")
	}

	result, err := s.aiRouter.RunFoodAnalysis(description, imageData, mimeType)
	if err != nil {
		return nil, ErrBadRequest(err.Error())
	}
	return result, nil
}

func (s *DietService) TrainerUpdateTargets(ctx context.Context, trainerUserID, clientID string, body map[string]interface{}) (map[string]interface{}, error) {
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}
	result, err := s.profileRepo.UpsertTargets(ctx, clientID, body)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return result, nil
}

func (s *DietService) TrainerGetMealsByDate(ctx context.Context, trainerUserID, clientID, date string) (interface{}, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}
	meals, err := s.mealRepo.FindByUserAndDate(ctx, clientID, date)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return meals, nil
}

func (s *DietService) TrainerGetStatus(ctx context.Context, trainerUserID, clientID, date string) (map[string]interface{}, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}
	return s.getDailyNutritionalStatus(ctx, clientID, date)
}

// ── Internal helpers ────────────────────────────────────────────────────────────

func (s *DietService) getDailyNutritionalStatus(ctx context.Context, userID, date string) (map[string]interface{}, error) {
	targets, err := s.profileRepo.GetTargets(ctx, userID)
	if err != nil {
		return nil, err
	}

	meals, err := s.mealRepo.FindByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	consumed := map[string]float64{"calories": 0, "protein": 0, "carbs": 0, "fat": 0, "fiber": 0}
	for _, meal := range meals {
		var items []map[string]interface{}
		json.Unmarshal(meal.Data, &items)
		for _, item := range items {
			consumed["calories"] += ToFloat(item["calories"])
			consumed["protein"] += ToFloat(item["protein"])
			consumed["carbs"] += ToFloat(item["carbs"])
			consumed["fat"] += ToFloat(item["fats"])
			consumed["fiber"] += ToFloat(item["fiber"])
		}
	}

	var targetCal, targetProt, targetCarbs, targetFat, targetFiber float64
	if targets != nil {
		targetCal = ToFloat(targets["targetCalories"])
		targetProt = ToFloat(targets["targetProtein"])
		targetCarbs = ToFloat(targets["targetCarbs"])
		targetFat = ToFloat(targets["targetFat"])
		targetFiber = ToFloat(targets["targetFiber"])
	}

	return map[string]interface{}{
		"targets": map[string]interface{}{
			"calories": targetCal,
			"protein":  targetProt,
			"carbs":    targetCarbs,
			"fat":      targetFat,
			"fiber":    targetFiber,
		},
		"consumed": map[string]interface{}{
			"calories": math.Round(consumed["calories"]),
			"protein":  math.Round(consumed["protein"]*10) / 10,
			"carbs":    math.Round(consumed["carbs"]*10) / 10,
			"fat":      math.Round(consumed["fat"]*10) / 10,
			"fiber":    math.Round(consumed["fiber"]*10) / 10,
		},
	}, nil
}

func (s *DietService) generateDietSuggestion(ctx context.Context, userID string) (map[string]interface{}, error) {
	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil || profile == nil {
		return nil, fmt.Errorf("user profile not found")
	}

	today := time.Now().Format("2006-01-02")
	status, err := s.getDailyNutritionalStatus(ctx, userID, today)
	if err != nil {
		return nil, err
	}

	targets := status["targets"].(map[string]interface{})
	consumed := status["consumed"].(map[string]interface{})

	remaining := map[string]float64{
		"calories": math.Max(0, ToFloat(targets["calories"])-ToFloat(consumed["calories"])),
		"protein":  math.Max(0, ToFloat(targets["protein"])-ToFloat(consumed["protein"])),
		"carbs":    math.Max(0, ToFloat(targets["carbs"])-ToFloat(consumed["carbs"])),
		"fat":      math.Max(0, ToFloat(targets["fat"])-ToFloat(consumed["fat"])),
		"fiber":    math.Max(0, ToFloat(targets["fiber"])-ToFloat(consumed["fiber"])),
	}

	mealsPerDay := 3
	if v := profile["mealsPerDay"]; v != nil {
		if mp, ok := v.(*int); ok && mp != nil {
			mealsPerDay = *mp
		}
	}

	userContext := fmt.Sprintf(`Profile: age %v, gender %v, height %vcm, weight %vkg, goal %v.
Meals per day: %d.

Daily Nutritional Status:
- Targets: { calories: %vkcal, protein: %vg, carbs: %vg, fat: %vg, fiber: %vg }
- Consumed Today: { calories: %vkcal, protein: %vg, carbs: %vg, fat: %vg, fiber: %vg }
- Gap (Lacking): { calories: %.0fkcal, protein: %.1fg, carbs: %.1fg, fat: %.1fg, fiber: %.1fg }`,
		profile["age"], profile["gender"], profile["heightCm"], profile["weightKg"], profile["goal"],
		mealsPerDay,
		targets["calories"], targets["protein"], targets["carbs"], targets["fat"], targets["fiber"],
		consumed["calories"], consumed["protein"], consumed["carbs"], consumed["fat"], consumed["fiber"],
		remaining["calories"], remaining["protein"], remaining["carbs"], remaining["fat"], remaining["fiber"],
	)

	aiResult, err := s.aiRouter.RunDietSuggestion(userContext)
	if err != nil {
		return nil, err
	}

	explanation, _ := aiResult["explanation"].(string)
	itemsJSON, _ := json.Marshal(aiResult["items"])

	return s.suggestionRepo.Upsert(ctx, userID, today, explanation, itemsJSON)
}
