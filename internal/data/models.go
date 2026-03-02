package data

import "database/sql"

type Models struct {
	Users UserModel
}

// NewModels returns a new Models struct with all the models wired to the provided db connection.
func NewModels(db *sql.DB) *Models {
	return &Models{
		Users: UserModel{
			DB: db,
		},
	}
}
