package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/heracle/pt.heracle.fit.go/internal/models"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, name, bio, avatar_url, email, google_access_token, created_at, updated_at
		 FROM users WHERE email = $1`, email).Scan(
		&u.ID, &u.Username, &u.Name, &u.Bio, &u.AvatarURL, &u.Email, &u.GoogleAccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, name, bio, avatar_url, email, google_access_token, created_at, updated_at
		 FROM users WHERE id = $1`, id).Scan(
		&u.ID, &u.Username, &u.Name, &u.Bio, &u.AvatarURL, &u.Email, &u.GoogleAccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, username, name, bio, avatar_url, email, google_access_token, created_at, updated_at
		 FROM users WHERE username = $1`, username).Scan(
		&u.ID, &u.Username, &u.Name, &u.Bio, &u.AvatarURL, &u.Email, &u.GoogleAccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

func (r *UserRepo) FindUsernamesStartingWith(ctx context.Context, prefix string) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT username FROM users WHERE username LIKE $1`, prefix+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		usernames = append(usernames, u)
	}
	return usernames, nil
}

func (r *UserRepo) Create(ctx context.Context, username, name, email string, avatarURL, googleAccessToken *string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (id, username, name, email, avatar_url, google_access_token, created_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW())
		 RETURNING id, username, name, bio, avatar_url, email, google_access_token, created_at, updated_at`,
		username, name, email, avatarURL, googleAccessToken).Scan(
		&u.ID, &u.Username, &u.Name, &u.Bio, &u.AvatarURL, &u.Email, &u.GoogleAccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	return &u, err
}

func (r *UserRepo) UpdateGoogleAccessToken(ctx context.Context, id string, token string) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx,
		`UPDATE users SET google_access_token = $2 WHERE id = $1
		 RETURNING id, username, name, bio, avatar_url, email, google_access_token, created_at, updated_at`,
		id, token).Scan(
		&u.ID, &u.Username, &u.Name, &u.Bio, &u.AvatarURL, &u.Email, &u.GoogleAccessToken, &u.CreatedAt, &u.UpdatedAt,
	)
	return &u, err
}

func (r *UserRepo) GetProfile(ctx context.Context, userID string) (map[string]interface{}, error) {
	var id, username, name, email string
	var avatarURL, bio *string
	var createdAt interface{}
	var updatedAt interface{}

	err := r.pool.QueryRow(ctx,
		`SELECT id, username, name, email, avatar_url, bio, created_at, updated_at
		 FROM users WHERE id = $1`, userID).Scan(
		&id, &username, &name, &email, &avatarURL, &bio, &createdAt, &updatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        id,
		"username":  username,
		"name":      name,
		"email":     email,
		"avatarUrl": avatarURL,
		"bio":       bio,
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}, nil
}
