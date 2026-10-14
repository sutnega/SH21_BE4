package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	Surname      string    `db:"surname"`
	Phone        string    `db:"phone"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}
