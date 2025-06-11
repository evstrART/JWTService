package repository

import (
	"JWTService/internal/models"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) SaveRefreshToken(ctx context.Context, userID int64, tokenID string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_id, expires_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, userID, tokenID, expiresAt)
	return err
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, tokenID string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `SELECT token_id, user_id, expires_at, revoked FROM refresh_tokens WHERE token_id = $1`

	err := r.db.GetContext(ctx, &token, query, tokenID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &token, nil
}
func (r *TokenRepository) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	query := `UPDATE refresh_tokens SET revoked = TRUE WHERE token_id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	return err
}

func (r *TokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *TokenRepository) DeleteAllByUserID(ctx context.Context, userID int64) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
