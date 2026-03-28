package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

type TrainerService struct {
	trainerRepo *repository.TrainerRepo
	userRepo    *repository.UserRepo
	profileRepo *repository.UserProfileRepo
	mealRepo    *repository.MealRepo
}

func NewTrainerService(
	trainerRepo *repository.TrainerRepo,
	userRepo *repository.UserRepo,
	profileRepo *repository.UserProfileRepo,
	mealRepo *repository.MealRepo,
) *TrainerService {
	return &TrainerService{
		trainerRepo: trainerRepo, userRepo: userRepo,
		profileRepo: profileRepo, mealRepo: mealRepo,
	}
}

func (s *TrainerService) GetClients(ctx context.Context, trainerUserID string) ([]map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	trainer, err := s.trainerRepo.FindByUserID(ctx, trainerUserID)
	if err != nil || trainer == nil {
		return nil, ErrForbidden("Trainer record not found")
	}

	clients, err := s.trainerRepo.FindClientsByTrainer(ctx, trainer.ID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	clientIDs := make([]string, len(clients))
	for i, cl := range clients {
		clientIDs[i], _ = cl["id"].(string)
	}

	todayMeals, _ := s.mealRepo.FindByUsersAndDate(ctx, clientIDs, today)

	for _, cl := range clients {
		clientID, _ := cl["id"].(string)
		var consumed float64
		for _, meal := range todayMeals {
			if meal.UserID == clientID {
				var items []map[string]interface{}
				json.Unmarshal(meal.Data, &items)
				for _, item := range items {
					consumed += ToFloat(item["calories"])
				}
			}
		}

		targetCal := ToFloat(cl["targetCalories"])
		progress := 0.0
		if targetCal > 0 {
			progress = math.Min(1, consumed/targetCal)
		}
		cl["progress"] = math.Round(progress*100) / 100
		delete(cl, "targetCalories")
	}

	return clients, nil
}

func (s *TrainerService) GetClientDetails(ctx context.Context, trainerUserID, clientID string) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	trainer, err := s.trainerRepo.FindByUserID(ctx, trainerUserID)
	if err != nil || trainer == nil {
		return nil, ErrForbidden("Trainer record not found")
	}

	assignment, err := s.trainerRepo.FindClientAssignment(ctx, clientID)
	if err != nil || assignment == nil || assignment.TrainerID != trainer.ID {
		return nil, ErrForbidden("Client is not assigned to you")
	}

	clientInfo, _ := s.trainerRepo.GetClientDetails(ctx, clientID)
	if clientInfo == nil {
		return nil, ErrNotFound("Client not found")
	}

	profile, _ := s.profileRepo.FindByUserID(ctx, clientID)

	clientMeals, _ := s.mealRepo.FindByUserAndDate(ctx, clientID, today)
	var consumed float64
	for _, meal := range clientMeals {
		var items []map[string]interface{}
		json.Unmarshal(meal.Data, &items)
		for _, item := range items {
			consumed += ToFloat(item["calories"])
		}
	}

	var targetCal float64
	if profile != nil {
		targetCal = ToFloat(profile["targetCalories"])
	}
	progress := 0.0
	if targetCal > 0 {
		progress = math.Min(1, consumed/targetCal)
	}

	result := map[string]interface{}{
		"id":         clientInfo["id"],
		"name":       clientInfo["name"],
		"email":      clientInfo["email"],
		"avatarUrl":  clientInfo["avatarUrl"],
		"assignedAt": assignment.AssignedAt,
		"progress":   math.Round(progress*100) / 100,
	}

	if profile != nil {
		result["goal"] = profile["goal"]
		result["age"] = profile["age"]
		result["gender"] = profile["gender"]
		result["heightCm"] = profile["heightCm"]
		result["weightKg"] = profile["weightKg"]
		result["bodyType"] = profile["bodyType"]
		result["fitnessLevel"] = profile["fitnessLevel"]
		result["bmi"] = profile["bmi"]
		result["targetCalories"] = profile["targetCalories"]
		result["targetProtein"] = profile["targetProtein"]
		result["targetCarbs"] = profile["targetCarbs"]
		result["targetFat"] = profile["targetFat"]
		result["targetFiber"] = profile["targetFiber"]
		result["injuries"] = profile["injuries"]
		result["dietaryPreference"] = profile["dietaryPreference"]
		result["workoutFrequencyPerWeek"] = profile["workoutFrequencyPerWeek"]
		result["preferredWorkoutType"] = profile["preferredWorkoutType"]
	}

	return result, nil
}

func (s *TrainerService) AddClient(ctx context.Context, trainerUserID, email string) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	trainer, _ := s.trainerRepo.FindByUserID(ctx, trainerUserID)
	if trainer == nil {
		return nil, ErrForbidden("Trainer record not found")
	}

	clientUser, _ := s.userRepo.FindByEmail(ctx, email)
	if clientUser == nil {
		return nil, ErrNotFound("User with email " + email + " not found")
	}

	existingAssignment, _ := s.trainerRepo.FindClientAssignment(ctx, clientUser.ID)
	if existingAssignment != nil {
		if existingAssignment.TrainerID == trainer.ID {
			return nil, &AppError{Status: 409, Message: "User is already your client"}
		}
		return nil, &AppError{Status: 409, Message: "User is already assigned to another trainer"}
	}

	assignment, err := s.trainerRepo.AddClient(ctx, trainer.ID, clientUser.ID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	profile, _ := s.profileRepo.FindByUserID(ctx, clientUser.ID)
	clientMeals, _ := s.mealRepo.FindByUserAndDate(ctx, clientUser.ID, today)

	var consumed float64
	for _, meal := range clientMeals {
		var items []map[string]interface{}
		json.Unmarshal(meal.Data, &items)
		for _, item := range items {
			consumed += ToFloat(item["calories"])
		}
	}

	var targetCal float64
	var goal interface{}
	if profile != nil {
		targetCal = ToFloat(profile["targetCalories"])
		goal = profile["goal"]
	}
	progress := 0.0
	if targetCal > 0 {
		progress = math.Min(1, consumed/targetCal)
	}

	return map[string]interface{}{
		"id":         clientUser.ID,
		"name":       clientUser.Name,
		"email":      clientUser.Email,
		"avatarUrl":  clientUser.AvatarURL,
		"assignedAt": assignment.AssignedAt,
		"goal":       goal,
		"progress":   math.Round(progress*100) / 100,
	}, nil
}

func (s *TrainerService) RemoveClient(ctx context.Context, trainerUserID, clientID string) error {
	trainer, _ := s.trainerRepo.FindByUserID(ctx, trainerUserID)
	if trainer == nil {
		return ErrForbidden("Trainer record not found")
	}

	assignment, _ := s.trainerRepo.FindClientAssignment(ctx, clientID)
	if assignment == nil || assignment.TrainerID != trainer.ID {
		return ErrForbidden("User is not your client")
	}

	return s.trainerRepo.RemoveClient(ctx, clientID)
}

func (s *TrainerService) AdminAddTrainer(ctx context.Context, email string, specialization *string, experience *int) (map[string]interface{}, error) {
	user, _ := s.userRepo.FindByEmail(ctx, email)
	if user == nil {
		return nil, ErrNotFound("User with email " + email + " not found")
	}

	existing, _ := s.trainerRepo.FindByUserID(ctx, user.ID)
	if existing != nil {
		return nil, &AppError{Status: 409, Message: "User is already a trainer"}
	}

	result, err := s.trainerRepo.CreateTrainer(ctx, user.ID, specialization, experience)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	result["user"] = map[string]interface{}{
		"id": user.ID, "username": user.Username, "name": user.Name, "email": user.Email,
	}
	return result, nil
}
