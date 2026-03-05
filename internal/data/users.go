package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Username  string
	Email     string
	Password  password
	Activated bool // not in use at the time
	CreatedAt time.Time
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

// UserModel is used to access the user data layer and execute db operations.
type UserModel struct {
	DB *sql.DB
}

// Insert inserts a new user into the database or returns an appropriate error.
func (um *UserModel) Insert(ctx context.Context, u *User) error {
	query := `INSERT INTO users (id, name, username, email, password_hash, activated) 
		    VALUES ($1, $2, $3, $4, $5, $6)
		    RETURNING created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{u.ID, u.Name, u.Username, u.Email, u.Password.Hash, u.Activated}

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return userConstraintErrs(err, "couldn't insert user")
	}
	return nil
}

// GetByUsername returns the user from the database for the provided username
// or returns an ErrRecordNotFound.
func (um *UserModel) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, name, username, email, password_hash, activated, created_at, updated_at 
		    FROM users 
		    WHERE username = $1`
	var u User

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.Name,
		&u.Username,
		&u.Email, &u.Password.Hash,
		&u.Activated,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("couldn't get user: %w", err)
		}
	}
	return &u, nil
}

// GetByEmail returns the user from the database for the provided email or returns
// ErrRecordNotFound.
func (um *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, username, email, password_hash, activated, created_at, updated_at 
		    FROM users 
		    WHERE email = $1`
	var u User

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Name,
		&u.Username,
		&u.Email, &u.Password.Hash,
		&u.Activated,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}
	return &u, nil
}

// Update updates the user in the database or returns an appropriate error.
func (um *UserModel) Update(ctx context.Context, user *User) error {
	query := `UPDATE users SET name=$1, username=$2, email=$3, activated=$4, updated_at=NOW()
		    WHERE id=$5
		    RETURNING updated_at`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{user.Name, user.Username, user.Email, user.Activated, user.ID}

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return userConstraintErrs(err, "couldn't update user")
	}
	return nil
}

// GetForToken returns the user for the provided auth token and scope if the provided
// token hasn't expired yet. If a user doesn't exist for the token, ErrRecordNotFound
// is returned.
func (um *UserModel) GetForToken(ctx context.Context, sc Scope, plainToken string) (*User, error) {
	query := `SELECT u.id, u.name, u.username, u.email, u.password_hash, u.activated, u.created_at, u.updated_at 
		    FROM users AS u
		    INNER JOIN tokens AS t ON u.id = t.user_id
		    WHERE t.hash = $1 AND t.scope=$2 AND t.expires_at > $3`
	var u User

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	hash := sha256.Sum256([]byte(plainToken))

	err := um.DB.QueryRowContext(ctx, query, hash[:], sc, time.Now()).Scan(
		&u.ID, &u.Name, &u.Username,
		&u.Email, &u.Password.Hash,
		&u.Activated,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to get user for token=%s: err=%w", plainToken, err)
		}
	}
	return &u, nil
}

// userConstraintErrs is a helper function to identify what type error occurred
// while inserting the user.
func userConstraintErrs(err error, msg string) error {
	if pqErr, ok := errors.AsType[*pq.Error](err); ok && pqErr.Code == "23505" {
		switch pqErr.Constraint {
		case "users_email_key":
			return ErrDuplicateEmail
		case "users_username_key":
			return ErrDuplicateUsername
		}
	}
	return fmt.Errorf("%s: %w", msg, err)
}
