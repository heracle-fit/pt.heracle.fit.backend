package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/heracle/pt.heracle.fit.go/internal/config"
	"github.com/heracle/pt.heracle.fit.go/internal/middleware"
	"github.com/heracle/pt.heracle.fit.go/internal/repository"
)

type AuthService struct {
	cfg         *config.Config
	userRepo    *repository.UserRepo
	trainerRepo *repository.TrainerRepo
	fireClient  *auth.Client
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepo, trainerRepo *repository.TrainerRepo) *AuthService {
	s := &AuthService{cfg: cfg, userRepo: userRepo, trainerRepo: trainerRepo}

	if cfg.FirebaseProjectID != "" && cfg.FirebaseClientEmail != "" && cfg.FirebasePrivateKey != "" {
		privateKey := strings.ReplaceAll(cfg.FirebasePrivateKey, "\\n", "\n")
		configMap := map[string]string{
			"type":         "service_account",
			"project_id":   cfg.FirebaseProjectID,
			"client_email": cfg.FirebaseClientEmail,
			"private_key":  privateKey,
			"token_uri":    "https://oauth2.googleapis.com/token",
		}

		credJSON, err := json.Marshal(configMap)
		if err != nil {
			log.Printf("⚠️  Firebase config JSON error: %v", err)
		} else {
			app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(credJSON))
			if err != nil {
				log.Printf("⚠️  Firebase init error: %v", err)
			} else {
				client, err := app.Auth(context.Background())
				if err != nil {
					log.Printf("⚠️  Firebase Auth client error: %v", err)
				} else {
					s.fireClient = client
					log.Println("✅ Firebase Admin SDK initialized")
				}
			}
		}
	}

	return s
}

type AuthResult struct {
	User  map[string]interface{}
	Token string
}

func (s *AuthService) AuthenticateGoogleToken(ctx context.Context, idToken, accessToken string) (*AuthResult, error) {
	log.Printf("AuthenticateGoogleToken called with idToken length: %d", len(idToken))
	if idToken == "" {
		return nil, fmt.Errorf("missing idToken")
	}
	if s.fireClient == nil {
		log.Println("❌ Firebase client is nil in AuthenticateGoogleToken")
		return nil, fmt.Errorf("firebase not initialized")
	}

	decoded, err := s.fireClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("❌ VerifyIDToken error: %v", err)
		return nil, fmt.Errorf("invalid Firebase ID token: %v", err)
	}
	log.Println("✅ VerifyIDToken success")

	email, _ := decoded.Claims["email"].(string)
	if email == "" {
		return nil, fmt.Errorf("firebase token does not contain an email")
	}

	name, _ := decoded.Claims["name"].(string)
	if name == "" {
		name = strings.Split(email, "@")[0]
	}
	picture, _ := decoded.Claims["picture"].(string)

	var at *string
	if accessToken != "" {
		at = &accessToken
	}
	var pic *string
	if picture != "" {
		pic = &picture
	}

	return s.findOrCreateUser(ctx, email, name, pic, at)
}

func (s *AuthService) AdminLogin(ctx context.Context, username, password string) (*AuthResult, error) {
	if username != s.cfg.AdminUsername || password != s.cfg.AdminPassword {
		return nil, fmt.Errorf("invalid admin credentials")
	}

	token, err := middleware.GenerateJWT(s.cfg.JWTSecret, "admin-id", "", "admin", s.cfg.AdminUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &AuthResult{
		User:  map[string]interface{}{"id": "admin-id", "username": s.cfg.AdminUsername, "role": "admin"},
		Token: token,
	}, nil
}

func (s *AuthService) GetDevToken(ctx context.Context, email string) (*AuthResult, error) {
	if email == "" {
		email = "sanjaysagar.main@gmail.com"
	}
	return s.findOrCreateUser(ctx, email, strings.Split(email, "@")[0], nil, nil)
}

func (s *AuthService) findOrCreateUser(ctx context.Context, email, displayName string, avatarURL, accessToken *string) (*AuthResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("database error")
	}

	if user == nil {
		baseUsername := strings.Split(email, "@")[0]
		baseUsername = sanitizeUsername(baseUsername)
		if len(baseUsername) > 30 {
			baseUsername = baseUsername[:30]
		}
		if baseUsername == "" {
			baseUsername = "user"
		}

		existingUsernames, _ := s.userRepo.FindUsernamesStartingWith(ctx, baseUsername)
		usernameSet := make(map[string]bool)
		for _, u := range existingUsernames {
			usernameSet[u] = true
		}

		username := baseUsername
		suffix := 0
		for usernameSet[username] {
			suffix++
			username = fmt.Sprintf("%s%d", baseUsername, suffix)
		}

		user, err = s.userRepo.Create(ctx, username, displayName, email, avatarURL, accessToken)
		if err != nil {
			return nil, fmt.Errorf("failed to create user")
		}
	} else if accessToken != nil && *accessToken != "" {
		if user.GoogleAccessToken == nil || *user.GoogleAccessToken != *accessToken {
			user, err = s.userRepo.UpdateGoogleAccessToken(ctx, user.ID, *accessToken)
			if err != nil {
				log.Printf("Failed to update access token: %v", err)
			}
		}
	}

	role := "user"
	trainer, _ := s.trainerRepo.FindByUserID(ctx, user.ID)
	if trainer != nil {
		role = "trainer"
	}

	token, err := middleware.GenerateJWT(s.cfg.JWTSecret, user.ID, user.Email, role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}

	return &AuthResult{
		User: map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"name":      user.Name,
			"email":     user.Email,
			"avatarUrl": user.AvatarURL,
			"bio":       user.Bio,
			"role":      role,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
		Token: token,
	}, nil
}

func sanitizeUsername(s string) string {
	var result strings.Builder
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '.' || c == '_' || c == '-' {
			result.WriteRune(c)
		}
	}
	return result.String()
}
