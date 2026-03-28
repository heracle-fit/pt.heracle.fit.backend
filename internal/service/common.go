package service

import (
	"context"
	"fmt"

	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

// AppError represents a service-layer error with an HTTP status code.
type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

// VerifyTrainerClient checks that the trainer is assigned to the given client.
func VerifyTrainerClient(ctx context.Context, trainerRepo *repository.TrainerRepo, trainerUserID, clientID string) *AppError {
	trainer, err := trainerRepo.FindByUserID(ctx, trainerUserID)
	if err != nil || trainer == nil {
		return &AppError{Status: 403, Message: "Trainer record not found for this user"}
	}

	assignment, err := trainerRepo.FindClientAssignment(ctx, clientID)
	if err != nil || assignment == nil || assignment.TrainerID != trainer.ID {
		return &AppError{Status: 403, Message: "You are not assigned to this client"}
	}
	return nil
}

// ToFloat converts various numeric types to float64.
func ToFloat(v interface{}) float64 {
	switch vt := v.(type) {
	case float64:
		return vt
	case int:
		return float64(vt)
	case *float64:
		if vt != nil {
			return *vt
		}
	case *int:
		if vt != nil {
			return float64(*vt)
		}
	}
	return 0
}

// ErrNotFound creates a 404 AppError.
func ErrNotFound(msg string) *AppError {
	return &AppError{Status: 404, Message: msg}
}

// ErrForbidden creates a 403 AppError.
func ErrForbidden(msg string) *AppError {
	return &AppError{Status: 403, Message: msg}
}

// ErrBadRequest creates a 400 AppError.
func ErrBadRequest(msg string) *AppError {
	return &AppError{Status: 400, Message: msg}
}

// ErrInternal creates a 500 AppError.
func ErrInternal(msg string) *AppError {
	return &AppError{Status: 500, Message: msg}
}

// Wrap wraps a Go error into a 500 AppError.
func Wrap(err error) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{Status: 500, Message: fmt.Sprintf("internal error: %v", err)}
}
