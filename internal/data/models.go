package data

import "database/sql"

type Models struct {
	Collections CollectionModel
	Users       UserModel
	Tokens      TokenModel
}

// NewModels returns a new Models struct with all the models wired to the provided db connection.
func NewModels(db *sql.DB) *Models {
	return &Models{
		Collections: CollectionModel{
			DB: db,
		},
		Users: UserModel{
			DB: db,
		},
		Tokens: TokenModel{
			DB: db,
		},
	}
}
