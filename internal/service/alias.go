package service

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"url-short/internal/domain"
)

func AliasValidate(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("%w: alias is null", domain.ErrValidation)
	}

	if utf8.RuneCountInString(s) > 20 {
		return fmt.Errorf("%w: alias too big", domain.ErrValidation)
	}

	if utf8.RuneCountInString(s) < 5 {
		return fmt.Errorf("%w: alias too short", domain.ErrValidation)
	}

	for _, r := range s {
		if !strings.ContainsRune(domain.AliasAlphabet, r) {
			return fmt.Errorf("%w: alias содержит недопустимый символ: %c", domain.ErrValidation, r)
		}
	}

	return nil
}
