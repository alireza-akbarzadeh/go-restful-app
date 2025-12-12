package validation

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		// Register custom validators
		err := validate.RegisterValidation("slug", validateSlug)
		if err != nil {
			return
		}
		err = validate.RegisterValidation("strong_password", validateStrongPassword)
		if err != nil {
			return
		}

		// Register custom tag name func
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// validateSlug validates URL-friendly slug format
func validateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	if slug == "" {
		return false
	}

	// Check if slug contains only lowercase letters, numbers, and hyphens
	for _, r := range slug {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}

	// Slug cannot start or end with hyphen
	return !strings.HasPrefix(slug, "-") && !strings.HasSuffix(slug, "-")
}

// validateStrongPassword validates password strength
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Error represents a structured validation error
type Error struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// FormatValidationErrors formats validator errors into a user-friendly format
func FormatValidationErrors(err error) []Error {
	var errors []Error

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			field := err.Field()
			value := err.Value().(string)
			message := getValidationMessage(err.Tag(), err.Param())

			errors = append(errors, Error{
				Field:   field,
				Value:   value,
				Message: message,
			})
		}
	}

	return errors
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email address"
	case "min":
		return "Must be at least " + param + " characters"
	case "max":
		return "Must be no more than " + param + " characters"
	case "slug":
		return "Must be a valid URL slug (lowercase letters, numbers, hyphens only)"
	case "strong_password":
		return "Password must contain at least one uppercase letter, lowercase letter, number, and special character"
	default:
		return "Invalid value"
	}
}
