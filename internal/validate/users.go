package validate

import (
	"Continuity/internal/data"
	"regexp"
)

var EmailRegex = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9]" +
		"(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?" +
		"(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

// ValidateEmail checks if the given email is valid.
func ValidateEmail(v *Validator, email string) {
	v.CheckField(NotBlank(email), "email", "must not be blank")
	v.CheckField(Matches(email, EmailRegex), "email", "invalid email address")
}

// ValidatePlainPassword checks if the given plaintext password is valid.
func ValidatePlainPassword(v *Validator, password string) {
	v.CheckField(NotBlank(password), "password", "must not be blank")
	v.CheckField(MinChars(password, 8), "password", "must be greater than 8 characters")
	v.CheckField(MaxChars(password, 72), "password", "must be lesser than 72 characters")
}

// ValidateUser checks if the given user is valid.
func ValidateUser(v *Validator, u *data.User) {
	v.CheckField(NotBlank(u.Name), "name", "must not be blank")
	v.CheckField(MaxChars(u.Name, 128), "name", "must be lesser than 128 characters")

	v.CheckField(NotBlank(u.Username), "username", "must not be blank")
	v.CheckField(MaxChars(u.Username, 128), "username", "must be lesser than 128 characters")

	ValidateEmail(v, u.Email)
	ValidatePlainPassword(v, u.Password.PlainText)

	if u.Password.Hash == nil {
		panic("missing password hash for user")
	}
}

// ValidatePlainToken checks if the given plaintext token is valid.
func ValidatePlainToken(v *Validator, plainToken string) {
	v.Check(NotBlank(plainToken), "plainToken", "must not be blank")
	v.Check(len(plainToken) == 26, "plainToken", "must be 26 characters long")
}
