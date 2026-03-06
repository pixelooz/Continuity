package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/google/uuid"
)

type Scope string

const (
	ScopeActivation     Scope = "activation"
	ScopeAuthentication Scope = "authentication"
)

// Token holds the data necessary for a successful authentication.
type Token struct {
	UserID    uuid.UUID
	PlainText string
	Hash      []byte
	Scope     Scope
	ExpiresAt time.Time
}

type TokenModel struct {
	DB *sql.DB
}

// generateToken generates a new session token and initializes a new Token struct
// with the provided data.
func generateToken(userID uuid.UUID, ttl time.Duration, scope Scope) *Token {
	ss := &Token{
		UserID:    userID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(ttl),
	}
	random := make([]byte, 16)
	_, _ = rand.Read(random) // method never returns an error, so don't handle; read doc.

	ss.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(random)
	hash := sha256.Sum256([]byte(ss.PlainText))
	ss.Hash = hash[:]
	return ss
}

// NewToken generates a new session token and inserts it into the database.
func (tm *TokenModel) NewToken(ctx context.Context, userID uuid.UUID, ttl time.Duration, sc Scope) (*Token, error) {
	token := generateToken(userID, ttl, sc)
	if err := tm.Insert(ctx, token); err != nil {
		return nil, err
	}
	return token, nil
}

// Insert inserts a new token into the database.
func (tm *TokenModel) Insert(ctx context.Context, token *Token) error {
	query := `INSERT INTO tokens(user_id, hash, scope, expires_at) VALUES ($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	args := []any{token.UserID, token.Hash, token.Scope, token.ExpiresAt}

	_, err := tm.DB.ExecContext(ctx, query, args...)
	return err
}

// DeleteAllForUser deletes all tokens for the provided UserID and scope.
func (tm *TokenModel) DeleteAllForUser(ctx context.Context, userID uuid.UUID, sc Scope) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := tm.DB.ExecContext(ctx, query, userID, sc)
	return err
}
