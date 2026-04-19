package rules

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func StringMinLength(value string, min int, attribute string, errs *[]ValidationError) {
	if len(value) < min {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s length must be at least %d", attribute, min),
		})
	}
}

func StringMaxLength(value string, max int, attribute string, errs *[]ValidationError) {
	if len(value) > max {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s length must be less than %d", attribute, max),
		})
	}
}

func IntMinLength(value int, min int, attribute string, errs *[]ValidationError) {
	if value < min {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s must be at least %d", attribute, min),
		})
	}
}

func IntMaxLength(value int, max int, attribute string, errs *[]ValidationError) {
	if value > max {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s must be less than %d", attribute, max),
		})
	}
}

func MustContainsAny(value string, characters string, number int, attribute string, errs *[]ValidationError) {
	count := 0
	for _, char := range value {
		if strings.ContainsRune(characters, char) {
			count++
			if count >= number {
				return
			}
		}
	}
	*errs = append(*errs, ValidationError{
		Field:   attribute,
		Message: fmt.Sprintf("%s must contain at least %d of these chars (%s)", attribute, number, characters),
	})
}

func MustNotContainsAny(value string, characters string, attribute string, errs *[]ValidationError) {
	if strings.ContainsAny(value, characters) {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s must not contain any of these chars (%s)", attribute, characters),
		})
	}
}

func MustContains(value string, word string, attribute string, errs *[]ValidationError) {
	if !strings.Contains(value, word) {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s must contain the word (%s)", attribute, word),
		})
	}
}

func MustNotContains(value string, word string, attribute string, errs *[]ValidationError) {
	if strings.Contains(value, word) {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s must not contain the forbidden word (%s)", attribute, word),
		})
	}
}

func StringStart(value string, prefix string, attribute string, errs *[]ValidationError) {
	if !strings.HasPrefix(value, prefix) {
		*errs = append(*errs, ValidationError{
			Field:   attribute,
			Message: fmt.Sprintf("%s must be prefixed by %s", attribute, prefix),
		})
	}
}
