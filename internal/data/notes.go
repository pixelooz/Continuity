package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Note is a single user-created note.
type Note struct {
	ID           uuid.UUID
	UserId       uuid.UUID
	CollectionId uuid.UUID
	Title        string
	Content      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type NoteModel struct {
	DB *sql.DB
}

func (nm *NoteModel) Insert(ctx context.Context, note *Note) error {
	query := `INSERT INTO notes (id, user_id, collection_id, title, content) 
			  VALUES ($1, $2, $3, $4, $5)
			  RETURNING created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{note.ID, note.UserId, note.CollectionId, note.Title, note.Content}

	err := nm.DB.QueryRowContext(ctx, query, args...).Scan(&note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return noteConstraintErrs(err, "couldn't insert note")
	}
	return nil
}

func (nm *NoteModel) GetByID(ctx context.Context, noteID uuid.UUID) (*Note, error) {
	query := `SELECT id, user_id, collection_id, title, content, created_at, updated_at 
			  FROM notes WHERE id = $1`
	var note Note

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := nm.DB.QueryRowContext(ctx, query, noteID).Scan(
		&note.ID, &note.UserId,
		&note.CollectionId,
		&note.Title, &note.Content,
		&note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("couldn't get note: %w", err)
		}
	}
	return &note, nil
}

func noteConstraintErrs(err error, msg string) error {
	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		switch pqErr.Code {
		case "23514":
			if pqErr.Constraint == "notes_title_not_empty" {
				return ErrNoteNameEmpty
			}
			// todo: add for note body as well.
		case "23505":
			if pqErr.Constraint == "notes_unique_per_collection" {
				return ErrDuplicateNote
			}
		}
	}
	return fmt.Errorf("%s: %w", msg, err)
}
