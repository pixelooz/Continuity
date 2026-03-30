package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// Collection is a collection of notes or other sub-collections.
type Collection struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ParentID  uuid.NullUUID
	Name      string
	IsRoot    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CollectionModel struct {
	DB *sql.DB
}

// Insert inserts a new collection into the database.
func (cm *CollectionModel) Insert(ctx context.Context, colxn *Collection) error {
	query := `INSERT INTO collections(id, user_id, parent_id, name) 
		    VALUES ($1, $2, $3, $4)
		    RETURNING created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{colxn.ID, colxn.UserID, colxn.ParentID, colxn.Name}

	err := cm.DB.QueryRowContext(ctx, query, args...).Scan(&colxn.CreatedAt, &colxn.UpdatedAt)
	if err != nil {
		return collectionConstraintErrs(err, "couldn't insert collection")
	}
	return nil
}

// GetRootCollection returns the root collection.
func (cm *CollectionModel) GetRootCollection(ctx context.Context, userID uuid.UUID) (*Collection, error) {
	query := `SELECT id, user_id, parent_id, name, is_root, created_at, updated_at 
		    FROM collections 
		    WHERE user_id = $1 AND is_root = TRUE`
	var c Collection

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := cm.DB.QueryRowContext(ctx, query, userID).Scan(
		&c.ID, &c.UserID,
		&c.ParentID,
		&c.Name,
		&c.IsRoot,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("couldn't get collection: %w", err)
		}
	}
	return &c, nil
}

// GetAllForParentId returns a slice of collections for the given parent id.
// todo: this method will probably have a join query also returning the pages for this parent id
func (cm *CollectionModel) GetAllForParentId(ctx context.Context, pID uuid.UUID) ([]*Collection, error) {
	query := `SELECT id, user_id, parent_id, name, is_root, created_at, updated_at 
		    FROM collections 
		    WHERE parent_id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := cm.DB.QueryContext(ctx, query, pID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rowErr := rows.Close(); rowErr != nil {
			log.Err(rowErr).Msg("couldn't close rows")
		}
	}()
	var collections []*Collection

	for rows.Next() {
		var c Collection
		err = rows.Scan(&c.ID, &c.UserID,
			&c.ParentID, &c.Name,
			&c.IsRoot,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		collections = append(collections, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("during scanning: %w", err)
	}
	return collections, nil
}

// GetById returns the collection from the database for the provided id or returns
// ErrRecordNotFound.
func (cm *CollectionModel) GetById(ctx context.Context, colxnID uuid.UUID) (*Collection, error) {
	query := `SELECT id, user_id, parent_id ,name, is_root, created_at, updated_at 
		    FROM collections WHERE id = $1`
	var colxn Collection

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := cm.DB.QueryRowContext(ctx, query, colxnID).Scan(
		&colxn.ID, &colxn.UserID,
		&colxn.ParentID,
		&colxn.Name,
		&colxn.IsRoot,
		&colxn.CreatedAt, &colxn.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("couldn't get collection: %w", err)
		}
	}
	return &colxn, nil
}

// Update updates the name or parent of the collection by id.
func (cm *CollectionModel) Update(ctx context.Context, colxn *Collection) error {
	query := `UPDATE collections SET name=$1, parent_id=$2, updated_at=NOW()
		    WHERE id=$3
		    RETURNING updated_at`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{colxn.Name, colxn.ParentID, colxn.ID}

	err := cm.DB.QueryRowContext(ctx, query, args...).Scan(&colxn.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRecordNotFound
		}
		return collectionConstraintErrs(err, "couldn't update colxn")
	}
	return nil
}

// Delete deletes the collection by the provided collection id and user id.
func (cm *CollectionModel) Delete(ctx context.Context, colxnID, userID uuid.UUID) error {
	query := `DELETE FROM collections WHERE id = $1 AND user_id = $2 AND is_root = FALSE`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := cm.DB.ExecContext(ctx, query, colxnID, userID)
	if err != nil {
		return collectionConstraintErrs(err, "couldn't delete collection")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// collectionConstraintErrs is a helper function to identify what type error occurred
// while inserting the collection.
func collectionConstraintErrs(err error, msg string) error {
	if pqErr, ok := errors.AsType[*pq.Error](err); ok {
		switch pqErr.Code {
		case "23514": // check_violation
			if pqErr.Constraint == "collections_name_not_empty" {
				return ErrCollectionNameEmpty
			}
		case "23503":
			if pqErr.Constraint == "notes_collection_id_fkey" {
				return ErrCollectionNotEmpty
			}
		case "23505": // unique_violation
			if pqErr.Constraint == "collection_unique_per_parent" {
				return ErrDuplicateCollection
			}
		}
	}
	return fmt.Errorf("%s: %w", msg, err)
}
