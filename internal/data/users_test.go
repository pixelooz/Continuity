package data

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestMapUserInsertErr(t *testing.T) {
	pqErr := &pq.Error{
		Code:       "23505",
		Constraint: "users_username_key",
	}
	err := userConstraintErrs(pqErr, "")
	require.Equal(t, err, ErrDuplicateUsername)
	require.NotEqual(t, err, ErrDuplicateEmail)

	pqErr = &pq.Error{
		Code:       "23505",
		Constraint: "users_email_key",
	}
	err = userConstraintErrs(pqErr, "")
	require.Equal(t, err, ErrDuplicateEmail)
	require.NotEqual(t, err, ErrDuplicateUsername)

	pqErr = &pq.Error{
		Code:       "22001",
		Constraint: "character too long",
	}
	err = userConstraintErrs(pqErr, "")
	require.NotEqual(t, err, ErrDuplicateEmail)
	require.NotEqual(t, err, ErrDuplicateUsername)
}
