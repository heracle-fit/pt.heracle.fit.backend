package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/heracle/pt.heracle.fit.go/internal/ai"
	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

type SleepService struct {
	sleepRepo *repository.SleepRepo
	aiRouter  *ai.AIRouter
}

func NewSleepService(sleepRepo *repository.SleepRepo, aiRouter *ai.AIRouter) *SleepService {
	return &SleepService{sleepRepo: sleepRepo, aiRouter: aiRouter}
}

func (s *SleepService) AddSleepData(ctx context.Context, userID string, newEntry map[string]interface{}) (interface{}, error) {
	existing, _ := s.sleepRepo.FindByUser(ctx, userID)

	var currentData []interface{}
	if existing != nil {
		if sd, ok := existing["sleepData"]; ok {
			if raw, ok := sd.(json.RawMessage); ok {
				json.Unmarshal(raw, &currentData)
			}
		}
	}

	currentData = append(currentData, newEntry)
	if len(currentData) > 7 {
		currentData = currentData[len(currentData)-7:]
	}

	updatedJSON, _ := json.Marshal(currentData)

	if existing != nil {
		id, _ := existing["id"].(string)
		if err := s.sleepRepo.UpdateSleepData(ctx, id, updatedJSON); err != nil {
			return nil, fmt.Errorf("database error")
		}
		existing["sleepData"] = json.RawMessage(updatedJSON)
		return existing, nil
	}

	newID, err := s.sleepRepo.Create(ctx, userID, updatedJSON)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return map[string]interface{}{
		"id": newID, "userId": userID, "sleepData": json.RawMessage(updatedJSON),
	}, nil
}

func (s *SleepService) GetSleepData(ctx context.Context, userID string) (interface{}, error) {
	existing, _ := s.sleepRepo.FindByUser(ctx, userID)
	if existing == nil {
		return map[string]interface{}{"sleepData": []interface{}{}}, nil
	}

	var currentData []interface{}
	if sd, ok := existing["sleepData"]; ok {
		if raw, ok := sd.(json.RawMessage); ok {
			json.Unmarshal(raw, &currentData)
		}
	}

	if len(currentData) > 7 {
		currentData = currentData[len(currentData)-7:]
		updatedJSON, _ := json.Marshal(currentData)
		id, _ := existing["id"].(string)
		s.sleepRepo.UpdateSleepData(ctx, id, updatedJSON)
	}

	return map[string]interface{}{"sleepData": currentData}, nil
}

func (s *SleepService) GetAIInsight(ctx context.Context, userID string) (map[string]interface{}, error) {
	today := time.Now().Format("2006-01-02")

	existing, _ := s.sleepRepo.FindByUser(ctx, userID)
	if existing != nil {
		if insight, ok := existing["insight"].(*string); ok && insight != nil && *insight != "" {
			if insightDate, ok := existing["insightDate"].(*string); ok && insightDate != nil && *insightDate == today {
				return map[string]interface{}{"insight": *insight}, nil
			}
		}
	}

	var sleepData []interface{}
	if existing != nil {
		if sd, ok := existing["sleepData"].(json.RawMessage); ok {
			json.Unmarshal(sd, &sleepData)
		}
	}

	sleepJSON, _ := json.Marshal(sleepData)
	result, err := s.aiRouter.RunSleepInsight(string(sleepJSON))
	if err != nil {
		log.Printf("[SleepService] AI Insight failed: %v", err)
		return nil, fmt.Errorf("failed to generate sleep insight")
	}

	insightText, _ := result["insight"].(string)

	if existing != nil {
		id, _ := existing["id"].(string)
		s.sleepRepo.UpdateInsight(ctx, id, insightText, today)
	} else {
		emptyData, _ := json.Marshal([]interface{}{})
		s.sleepRepo.CreateWithInsight(ctx, userID, emptyData, insightText, today)
	}

	return map[string]interface{}{"insight": insightText}, nil
}
