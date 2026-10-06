package service

import (
	"fmt"
	"strings"
	"unicode/utf8"
	"url-short/internal/domain"
)

func AliasValidate(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("Алиас не может быть пустым!")
	}

	if utf8.RuneCountInString(s) > 20 {
		return fmt.Errorf("Не может быть больше 20 символов")
	}

	if utf8.RuneCountInString(s) < 5 {
		return fmt.Errorf("Не может быть меньше 5 символов")
	}

	for _, r := range s {
		if !strings.ContainsRune(domain.AliasAlphabet, r) {
			return fmt.Errorf("alias содержит недопустимый символ: %c", r)
		}
	}

	if len(s) == 0 {
		return fmt.Errorf("Алиас не может быть пустым!")
	}

	return nil
}
