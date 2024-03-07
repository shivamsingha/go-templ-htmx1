package util

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
		return false
	}

	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+{}\[\]:;<>,.?~\\/\-=|"']`, password); !matched {
		return false
	}

	return true
}
