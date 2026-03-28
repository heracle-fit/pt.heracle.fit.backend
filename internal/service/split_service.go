package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

type SplitService struct {
	splitRepo   *repository.SplitRepo
	trainerRepo *repository.TrainerRepo
}

func NewSplitService(splitRepo *repository.SplitRepo, trainerRepo *repository.TrainerRepo) *SplitService {
	return &SplitService{splitRepo: splitRepo, trainerRepo: trainerRepo}
}

func (s *SplitService) GetMySplit(ctx context.Context, userID string) (map[string]interface{}, error) {
	split, err := s.splitRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	if split == nil {
		return nil, ErrNotFound("Workout split not found")
	}
	return split, nil
}

func (s *SplitService) TrainerGetClientSplit(ctx context.Context, trainerUserID, clientID string) (map[string]interface{}, error) {
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}

	split, err := s.splitRepo.FindByUserID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	if split == nil {
		return nil, ErrNotFound("Workout split not found")
	}
	return split, nil
}

func (s *SplitService) UpsertClientSplit(ctx context.Context, trainerUserID, clientID string, splitData json.RawMessage) (map[string]interface{}, error) {
	trainer, _ := s.trainerRepo.FindByUserID(ctx, trainerUserID)
	if trainer == nil {
		return nil, ErrForbidden("Trainer record not found")
	}

	assignment, _ := s.trainerRepo.FindClientAssignment(ctx, clientID)
	if assignment == nil || assignment.TrainerID != trainer.ID {
		return nil, ErrForbidden("You are not assigned to this client")
	}

	result, err := s.splitRepo.Upsert(ctx, clientID, trainer.ID, splitData)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return result, nil
}
