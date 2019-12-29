package stdapi

import (
	"regexp"
)

const nameValidatorPattern = `^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]$`

var nameValidator = regexp.MustCompile(nameValidatorPattern)

func ValidateName(name string) error {
	if !nameValidator.MatchString(name) {
		return ValidationError("name", nameValidatorPattern)
	}

	return nil
}
