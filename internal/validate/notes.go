package validate

import (
	"Continuity/internal/data"
)

func ValidateCollection(v *Validator, c *data.Collection) {
	v.CheckField(NotBlank(c.Name), "name", "must not be empty")
	v.CheckField(len(c.Name) <= 100, "name", "must be less than 100 characters")
}
