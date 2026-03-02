package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	constraintUniqueUsername = `pq: duplicate key value violates unique constraint "users_username_key"`
	constraintUniqueEmail    = `pq: duplicate key value violates unique constraint "users_email_key"`
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

// UserModel is used to access the user data layer and execute db operations.
type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID        uuid.UUID
	Name      string
	Username  string
	Email     string
	Password  password
	Activated bool
	CreateAt  time.Time
	UpdatedAt time.Time
}

type password struct {
	plainText *string
	hash      []byte
}
