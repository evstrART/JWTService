package models

import "time"

type RefreshToken struct {
	TokenID   string    `db:"token_id"`
	UserID    int64     `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	Revoked   bool      `db:"revoked"`
}
