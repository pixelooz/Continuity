package data

import "database/sql"

type Models struct {
	Collections CollectionModel
	Note        NoteModel
	Users       UserModel
	Tokens      TokenModel
}

// NewModels returns a new Models struct with all the models wired to the provided db connection.
func NewModels(db *sql.DB) *Models {
	return &Models{
		Collections: CollectionModel{
			DB: db,
		},
		Note: NoteModel{
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
