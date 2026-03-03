package validate

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Validator holds validation error values in key/value pairs.
type Validator struct {
	FieldErrors map[string]string
	Errors      map[string]string
}

// NewValidator returns a new initialized Validator struct.
func NewValidator() *Validator {
	return &Validator{
		Errors:      make(map[string]string),
		FieldErrors: make(map[string]string),
	}
}

// Valid returns true if there are no errors.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0 && len(v.FieldErrors) == 0
}

// AddFieldError adds a new error to the FieldErrors map.
func (v *Validator) AddFieldError(key, value string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = value
	}
}

// AddError adds a new error to the Errors map.
func (v *Validator) AddError(key, value string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	if _, exists := v.Errors[key]; !exists {
		v.FieldErrors[key] = value
	}
}

// Check adds a new error to the Errors map if the given condition is false.
func (v *Validator) Check(valid bool, key, value string) {
	if !valid {
		v.AddError(key, value)
	}
}

// CheckField adds a new error to the FieldErrors map if the given condition is false.
func (v *Validator) CheckField(valid bool, key, value string) {
	if !valid {
		v.AddFieldError(key, value)
	}
}

// NotBlank returns true if the given string is not empty after trimming whitespace.
func NotBlank(str string) bool {
	return strings.TrimSpace(str) != ""
}

// Matches returns true if the given string matches the given regular expression.
func Matches(str string, rx *regexp.Regexp) bool {
	return rx.MatchString(str)
}

// MaxChars returns true if the provided string is shorter than the given size.
func MaxChars(str string, size int) bool {
	return utf8.RuneCountInString(str) <= size
}

// MinChars returns true if the provided string is longer than the given size.
func MinChars(str string, size int) bool {
	return utf8.RuneCountInString(str) >= size
}
