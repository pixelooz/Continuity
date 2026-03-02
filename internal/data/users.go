package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
	PlainText string
	Hash      []byte
}

// SetPass hashes the provided plainText password and stores it in p.
func (p *password) SetPass(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return fmt.Errorf("couldn't hash password: %w", err)
	}
	p.Hash, p.PlainText = hash, plainText
	return nil
}

// CheckPass compares the provided plainText password with the stored hash.
func (p *password) CheckPass(plainText string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainText)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, fmt.Errorf("failed to check password: %w", err)
		}
	}
	return true, nil
}
