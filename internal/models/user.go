package models

import "time"

type User struct {
	Id            int64     `json:"id" db:"id"`
	Username      string    `db:"username" json:"username"`
	Email         string    `db:"email" json:"email"`
	Password_hash string    `db:"password_hash" json:"password_hash"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
