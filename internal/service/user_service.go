package service

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

type UserService struct {
	userRepo    *repository.UserRepo
	profileRepo *repository.UserProfileRepo
	trainerRepo *repository.TrainerRepo
}

func NewUserService(userRepo *repository.UserRepo, profileRepo *repository.UserProfileRepo, trainerRepo *repository.TrainerRepo) *UserService {
	return &UserService{userRepo: userRepo, profileRepo: profileRepo, trainerRepo: trainerRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (map[string]interface{}, error) {
	profile, err := s.userRepo.GetProfile(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return profile, nil
}

func (s *UserService) GetBodyMetrics(ctx context.Context, userID string) (map[string]interface{}, error) {
	metrics, err := s.profileRepo.GetBodyMetrics(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return metrics, nil
}

func (s *UserService) GetOnboardingStatus(ctx context.Context, userID string) (map[string]interface{}, error) {
	status, err := s.profileRepo.GetOnboardingStatus(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return status, nil
}

func (s *UserService) GetCalendarDetails(ctx context.Context, userID string) ([]byte, int, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, 404, fmt.Errorf("user not found")
	}
	if user.GoogleAccessToken == nil || *user.GoogleAccessToken == "" {
		return nil, 404, fmt.Errorf("Google Access Token not found")
	}

	req, _ := http.NewRequest("GET", "https://www.googleapis.com/calendar/v3/users/me/calendarList", nil)
	req.Header.Set("Authorization", "Bearer "+*user.GoogleAccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to fetch calendar details")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, 403, fmt.Errorf("Google Access Token is invalid or expired")
	}

	body, _ := io.ReadAll(resp.Body)
	return body, 200, nil
}

func (s *UserService) SaveBodyMetrics(ctx context.Context, userID string, body map[string]interface{}) (map[string]interface{}, error) {
	body = s.calculateDerivedMetrics(ctx, userID, body)
	result, err := s.profileRepo.UpsertBodyMetrics(ctx, userID, body)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err.Error())
	}
	return result, nil
}

func (s *UserService) SaveTargets(ctx context.Context, userID string, body map[string]interface{}) (map[string]interface{}, error) {
	result, err := s.profileRepo.UpsertTargets(ctx, userID, body)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return result, nil
}

func (s *UserService) TrainerSaveBodyMetrics(ctx context.Context, trainerUserID, clientID string, body map[string]interface{}) (map[string]interface{}, error) {
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}
	body = s.calculateDerivedMetrics(ctx, clientID, body)
	result, err := s.profileRepo.UpsertBodyMetrics(ctx, clientID, body)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return result, nil
}

func (s *UserService) calculateDerivedMetrics(ctx context.Context, userID string, body map[string]interface{}) map[string]interface{} {
	existing, _ := s.profileRepo.FindByUserID(ctx, userID)

	getFloat := func(key string) *float64 {
		if v, ok := body[key]; ok && v != nil {
			switch vt := v.(type) {
			case float64:
				return &vt
			case int:
				f := float64(vt)
				return &f
			}
		}
		if existing != nil {
			if v, ok := existing[key]; ok && v != nil {
				switch vt := v.(type) {
				case *float64:
					return vt
				case float64:
					return &vt
				}
			}
		}
		return nil
	}

	getInt := func(key string) *int {
		if v, ok := body[key]; ok && v != nil {
			switch vt := v.(type) {
			case float64:
				i := int(vt)
				return &i
			case int:
				return &vt
			}
		}
		if existing != nil {
			if v, ok := existing[key]; ok && v != nil {
				switch vt := v.(type) {
				case *int:
					return vt
				case int:
					return &vt
				}
			}
		}
		return nil
	}

	getString := func(key string) *string {
		if v, ok := body[key]; ok && v != nil {
			if s, ok := v.(string); ok {
				return &s
			}
		}
		if existing != nil {
			if v, ok := existing[key]; ok && v != nil {
				if s, ok := v.(*string); ok {
					return s
				}
				if s, ok := v.(string); ok {
					return &s
				}
			}
		}
		return nil
	}

	age := getInt("age")
	gender := getString("gender")
	heightCm := getFloat("heightCm")
	weightKg := getFloat("weightKg")
	goal := getString("goal")

	if heightCm != nil && weightKg != nil {
		bmi := *weightKg / math.Pow(*heightCm/100, 2)
		bmiRound := math.Round(bmi*10) / 10
		body["bmi"] = bmiRound
	}

	if age != nil && gender != nil && heightCm != nil && weightKg != nil {
		bmr := 10**weightKg + 6.25**heightCm - 5*float64(*age)
		if strings.EqualFold(*gender, "male") {
			bmr += 5
		} else {
			bmr -= 161
		}
		maintenance := int(math.Round(bmr * 1.2))
		body["maintenanceCalories"] = maintenance

		targetCal := maintenance
		if goal != nil {
			switch *goal {
			case "weight_loss":
				targetCal = maintenance - 500
			case "muscle_gain":
				targetCal = maintenance + 300
			}
		}
		body["targetCalories"] = targetCal

		targetProtein := math.Round(*weightKg*2.0*10) / 10
		targetFat := math.Round(*weightKg*0.7*10) / 10
		proteinCal := targetProtein * 4
		fatCal := targetFat * 9
		remaining := math.Max(0, float64(targetCal)-proteinCal-fatCal)
		targetCarbs := math.Round(remaining/4*10) / 10
		targetFiber := math.Round(float64(targetCal)/1000*14*10) / 10

		body["targetProtein"] = targetProtein
		body["targetFat"] = targetFat
		body["targetCarbs"] = targetCarbs
		body["targetFiber"] = targetFiber
	}

	return body
}
