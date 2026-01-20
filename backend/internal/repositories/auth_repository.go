package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	db *sqlx.DB
}

type AuthRepositoryInterface interface {
	CreateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	CreateRefreshToken(ctx context.Context, token *entities.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
}

func NewAuthRepository(db *sqlx.DB) AuthRepositoryInterface {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `
		INSERT INTO users (email, password, name, role, is_active) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx, query,
		user.Email, user.Password, user.Name, user.Role, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `SELECT * FROM users WHERE email = $1 AND is_active = true`

	user := &entities.User{}
	if err := r.db.GetContext(ctx, user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *authRepository) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `SELECT * FROM users WHERE id = $1 AND is_active = true`

	user := &entities.User{}
	if err := r.db.GetContext(ctx, user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *authRepository) CreateRefreshToken(ctx context.Context, token *entities.RefreshToken) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, is_revoked) 
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, token.UserID, token.Token, token.ExpiresAt, token.IsRevoked)
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		return err
	}

	return nil
}

func (r *authRepository) GetRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `SELECT * FROM refresh_tokens WHERE token = $1 AND is_revoked = false AND expires_at > NOW()`

	refreshToken := &entities.RefreshToken{}
	if err := r.db.GetContext(ctx, refreshToken, query, token); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error getting refresh token: %v", err)
		return nil, err
	}

	return refreshToken, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `UPDATE refresh_tokens SET is_revoked = true, revoked_at = NOW() WHERE token = $1`

	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		return err
	}

	return nil
}

func (r *authRepository) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	query := `UPDATE refresh_tokens SET is_revoked = true, revoked_at = NOW() WHERE user_id = $1 AND is_revoked = false`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error revoking all user tokens: %v", err)
		return err
	}

	return nil
}
