package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

const ExerciseImageBaseURL = "https://r2.heracle.fit/exercises"

type WorkoutService struct {
	profileRepo    *repository.UserProfileRepo
	exerciseRepo   *repository.ExerciseRepo
	sessionRepo    *repository.SessionRepo
	workoutLogRepo *repository.WorkoutLogRepo
	trainerRepo    *repository.TrainerRepo
}

func NewWorkoutService(
	profileRepo *repository.UserProfileRepo,
	exerciseRepo *repository.ExerciseRepo,
	sessionRepo *repository.SessionRepo,
	workoutLogRepo *repository.WorkoutLogRepo,
	trainerRepo *repository.TrainerRepo,
) *WorkoutService {
	return &WorkoutService{
		profileRepo: profileRepo, exerciseRepo: exerciseRepo,
		sessionRepo: sessionRepo, workoutLogRepo: workoutLogRepo,
		trainerRepo: trainerRepo,
	}
}

func (s *WorkoutService) GetTodayWorkout(ctx context.Context, userID string) (interface{}, error) {
	prefs, err := s.profileRepo.GetWorkoutPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	if prefs == nil || prefs["fitnessLevel"] == nil || prefs["workoutFrequencyPerWeek"] == nil ||
		prefs["preferredWorkoutType"] == nil || prefs["preferredWorkoutTime"] == nil {
		return nil, nil
	}

	days, _ := prefs["availableDays"].([]string)
	if len(days) == 0 {
		return nil, nil
	}

	return map[string]interface{}{
		"title":     "Suggested Muscle",
		"highlight": "Bicep & Back",
		"subtext":   "Optional for hypertrophy",
		"duration":  45,
		"intensity": "hard",
		"session":   GetStaticSessions(),
	}, nil
}

func (s *WorkoutService) GetStaticSessions() interface{} {
	return GetStaticSessions()
}

func (s *WorkoutService) GetExercises(ctx context.Context) ([]map[string]interface{}, error) {
	exercises, err := s.exerciseRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	result := make([]map[string]interface{}, len(exercises))
	for i, ex := range exercises {
		result[i] = map[string]interface{}{
			"id":            ex.ID,
			"name":          ex.Name,
			"secondaryInfo": ex.SecondaryInfo,
			"exerciseType":  ex.ExerciseType,
			"image":         fmt.Sprintf("%s/%d.jpg", ExerciseImageBaseURL, ex.ID),
		}
	}
	return result, nil
}

func (s *WorkoutService) GetWorkoutPreferences(ctx context.Context, userID string) (map[string]interface{}, error) {
	prefs, err := s.profileRepo.GetWorkoutPreferences(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	return prefs, nil
}

func (s *WorkoutService) SaveWorkoutPreferences(ctx context.Context, userID string, body map[string]interface{}) (map[string]interface{}, error) {
	result, err := s.profileRepo.UpsertWorkoutPreferences(ctx, userID, body)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err.Error())
	}
	return result, nil
}

// ── Session CRUD ────────────────────────────────────────────────────────────────

type CreateSessionRequest struct {
	Name        string          `json:"name"`
	Category    json.RawMessage `json:"category"`
	SessionData json.RawMessage `json:"sessionData"`
}

func (s *WorkoutService) CreateSession(ctx context.Context, userID string, req CreateSessionRequest) (map[string]interface{}, error) {
	result, err := s.sessionRepo.Create(ctx, userID, req.Name, req.Category, req.SessionData)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) GetSession(ctx context.Context, userID string, id int) (map[string]interface{}, error) {
	result, err := s.sessionRepo.FindByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	if result == nil {
		return nil, ErrNotFound(fmt.Sprintf("Session with ID %d not found", id))
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) GetUserSessions(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	results, err := s.sessionRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	for _, r := range results {
		r["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	}
	return results, nil
}

func (s *WorkoutService) UpdateSession(ctx context.Context, userID string, id int, body map[string]interface{}) (map[string]interface{}, error) {
	existing, _ := s.sessionRepo.FindByIDAndUser(ctx, id, userID)
	if existing == nil {
		return nil, ErrNotFound(fmt.Sprintf("Session with ID %d not found", id))
	}

	data := map[string]interface{}{}
	if v, ok := body["name"]; ok {
		data["name"] = v
	}
	if v, ok := body["category"]; ok {
		b, _ := json.Marshal(v)
		data["category"] = json.RawMessage(b)
	}
	if v, ok := body["sessionData"]; ok {
		b, _ := json.Marshal(v)
		data["sessionData"] = json.RawMessage(b)
	}

	result, err := s.sessionRepo.Update(ctx, id, data)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) DeleteSession(ctx context.Context, userID string, id int) error {
	existing, _ := s.sessionRepo.FindByIDAndUser(ctx, id, userID)
	if existing == nil {
		return ErrNotFound(fmt.Sprintf("Session with ID %d not found", id))
	}
	return s.sessionRepo.Delete(ctx, id)
}

// ── Trainer Session operations ──────────────────────────────────────────────────

func (s *WorkoutService) TrainerUpdateSession(ctx context.Context, trainerUserID, clientID string, sessionID int, body map[string]interface{}) (map[string]interface{}, error) {
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}

	existing, _ := s.sessionRepo.FindByIDAndUser(ctx, sessionID, clientID)
	if existing == nil {
		return nil, ErrNotFound(fmt.Sprintf("Session %d not found for client %s", sessionID, clientID))
	}

	data := map[string]interface{}{}
	if v, ok := body["name"]; ok {
		data["name"] = v
	}
	if v, ok := body["category"]; ok {
		b, _ := json.Marshal(v)
		data["category"] = json.RawMessage(b)
	}
	if v, ok := body["sessionData"]; ok {
		b, _ := json.Marshal(v)
		data["sessionData"] = json.RawMessage(b)
	}

	result, err := s.sessionRepo.Update(ctx, sessionID, data)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) TrainerCreateSession(ctx context.Context, trainerUserID, clientID string, req CreateSessionRequest) (map[string]interface{}, error) {
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}

	result, err := s.sessionRepo.Create(ctx, clientID, req.Name, req.Category, req.SessionData)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) TrainerGetSessions(ctx context.Context, trainerUserID, clientID string) ([]map[string]interface{}, error) {
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return nil, err
	}

	results, err := s.sessionRepo.FindByUser(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	for _, r := range results {
		r["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	}
	return results, nil
}

func (s *WorkoutService) TrainerDeleteSession(ctx context.Context, trainerUserID, clientID string, sessionID int) error {
	if err := VerifyTrainerClient(ctx, s.trainerRepo, trainerUserID, clientID); err != nil {
		return err
	}

	existing, _ := s.sessionRepo.FindByIDAndUser(ctx, sessionID, clientID)
	if existing == nil {
		return ErrNotFound(fmt.Sprintf("Session %d not found for client %s", sessionID, clientID))
	}

	return s.sessionRepo.Delete(ctx, sessionID)
}

// ── WorkoutLog CRUD ─────────────────────────────────────────────────────────────

type CreateWorkoutLogRequest struct {
	SessionID *int            `json:"sessionId"`
	LogData   json.RawMessage `json:"logData"`
	Notes     *string         `json:"notes"`
}

func (s *WorkoutService) CreateWorkoutLog(ctx context.Context, userID string, req CreateWorkoutLogRequest) (map[string]interface{}, error) {
	result, err := s.workoutLogRepo.Create(ctx, userID, req.SessionID, req.LogData, req.Notes)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) GetWorkoutLog(ctx context.Context, userID string, id int) (map[string]interface{}, error) {
	result, err := s.workoutLogRepo.FindByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	if result == nil {
		return nil, ErrNotFound(fmt.Sprintf("Workout log with ID %d not found", id))
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) GetWorkoutLogs(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	results, err := s.workoutLogRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	for _, r := range results {
		r["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	}
	return results, nil
}

func (s *WorkoutService) UpdateWorkoutLog(ctx context.Context, userID string, id int, body map[string]interface{}) (map[string]interface{}, error) {
	existing, _ := s.workoutLogRepo.FindByIDAndUser(ctx, id, userID)
	if existing == nil {
		return nil, ErrNotFound(fmt.Sprintf("Workout log with ID %d not found", id))
	}

	data := map[string]interface{}{}
	if v, ok := body["sessionId"]; ok {
		if f, ok := v.(float64); ok {
			i := int(f)
			data["sessionId"] = &i
		}
	}
	if v, ok := body["logData"]; ok {
		b, _ := json.Marshal(v)
		data["logData"] = json.RawMessage(b)
	}
	if v, ok := body["notes"]; ok {
		if s, ok := v.(string); ok {
			data["notes"] = &s
		}
	}

	result, err := s.workoutLogRepo.Update(ctx, id, data)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

func (s *WorkoutService) DeleteWorkoutLog(ctx context.Context, userID string, id int) error {
	existing, _ := s.workoutLogRepo.FindByIDAndUser(ctx, id, userID)
	if existing == nil {
		return ErrNotFound(fmt.Sprintf("Workout log with ID %d not found", id))
	}
	return s.workoutLogRepo.Delete(ctx, id)
}

func (s *WorkoutService) TrainerAddLogReview(ctx context.Context, trainerUserID string, logID int, review string) (map[string]interface{}, error) {
	trainer, err := s.trainerRepo.FindByUserID(ctx, trainerUserID)
	if err != nil || trainer == nil {
		return nil, ErrForbidden("Trainer record not found")
	}

	logEntry, err := s.workoutLogRepo.FindByID(ctx, logID)
	if err != nil || logEntry == nil {
		return nil, ErrNotFound(fmt.Sprintf("Workout log with ID %d not found", logID))
	}

	logUserID, _ := logEntry["userId"].(string)
	assignment, err := s.trainerRepo.FindClientAssignment(ctx, logUserID)
	if err != nil || assignment == nil || assignment.TrainerID != trainer.ID {
		return nil, ErrForbidden("You are not assigned to the owner of this log")
	}

	result, updateErr := s.workoutLogRepo.UpdatePTReview(ctx, logID, review)
	if updateErr != nil {
		return nil, fmt.Errorf("database error")
	}
	result["exerciseImageBaseUrl"] = ExerciseImageBaseURL
	return result, nil
}

// ── Static sessions ─────────────────────────────────────────────────────────────

func GetStaticSessions() []map[string]interface{} {
	imgURL := "https://pub-7ec42550dbda4d5db5e62b8a86f5f595.r2.dev/exercises/Heracle.jpg"
	return []map[string]interface{}{
		{
			"id": "session-001", "title": "Chest & Triceps", "content": "Push day focused on upper body strength",
			"category": "Strength", "exercisesCount": 3, "position": 1,
			"exercises": []map[string]interface{}{
				{"id": "ex-001", "name": "Bench Press", "desc": "Chest, Triceps, Shoulders",
					"image": imgURL,
					"sets":  []map[string]interface{}{{"kg": 60, "reps": 10}, {"kg": 65, "reps": 8}, {"kg": 70, "reps": 6}}},
				{"id": "ex-002", "name": "Tricep Dips", "desc": "Triceps, Chest",
					"image": imgURL,
					"sets":  []map[string]interface{}{{"kg": 0, "reps": 12}, {"kg": 0, "reps": 10}, {"kg": 0, "reps": 10}}},
				{"id": "ex-003", "name": "Incline Dumbbell Press", "desc": "Upper Chest, Shoulders",
					"image": imgURL,
					"sets":  []map[string]interface{}{{"kg": 24, "reps": 10}, {"kg": 26, "reps": 8}, {"kg": 28, "reps": 6}}},
			},
		},
		{
			"id": "session-002", "title": "Back & Biceps", "content": "Pull day targeting back width and bicep peak",
			"category": "Strength", "exercisesCount": 3, "position": 2,
			"exercises": []map[string]interface{}{
				{"id": "ex-004", "name": "Deadlift", "desc": "Lower Back, Glutes, Hamstrings",
					"image": imgURL,
					"sets":  []map[string]interface{}{{"kg": 80, "reps": 8}, {"kg": 90, "reps": 6}, {"kg": 100, "reps": 4}}},
				{"id": "ex-005", "name": "Pull-Ups", "desc": "Lats, Biceps, Rear Delts",
					"image": imgURL,
					"sets":  []map[string]interface{}{{"kg": 0, "reps": 10}, {"kg": 0, "reps": 8}, {"kg": 0, "reps": 8}}},
				{"id": "ex-006", "name": "Barbell Curl", "desc": "Biceps, Forearms",
					"image": imgURL,
					"sets":  []map[string]interface{}{{"kg": 30, "reps": 12}, {"kg": 35, "reps": 10}, {"kg": 35, "reps": 8}}},
			},
		},
	}
}

// ParseID parses a string ID to int, returning 0 on error.
func ParseID(s string) int {
	id, _ := strconv.Atoi(s)
	return id
}
