package validate

import (
	"Continuity/internal/data"
)

func ValidateCollection(v *Validator, c *data.Collection) {
	v.CheckField(NotBlank(c.Name), "name", "must not be empty")
	v.CheckField(len(c.Name) <= 100, "name", "must be less than 100 characters")
}

func ValidateNote(v *Validator, n *data.Note) {
	v.CheckField(NotBlank(n.Title), "name", "must not be empty")
	v.CheckField(len(n.Title) <= 100, "name", "must be less than 100 characters")
}
