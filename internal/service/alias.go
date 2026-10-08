package service

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/url"
	"strings"
	"unicode/utf8"
	"url-short/internal/domain"
)

// GenerateAlias случайный alias из AliasAlphabet длиной 8 символов.
func GenerateAlias() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = domain.AliasAlphabet[rand.IntN(len(domain.AliasAlphabet))]
	}
	return string(b)
}

func AliasValidate(s string) error {
	n := utf8.RuneCountInString(s)

	if n > 20 {
		return fmt.Errorf("%w: alias too long", domain.ErrValidation)
	}
	if n < 5 {
		return fmt.Errorf("%w: alias too short", domain.ErrValidation)
	}

	for _, r := range s {
		if !strings.ContainsRune(domain.AliasAlphabet, r) {
			return fmt.Errorf("%w: alias contains invalid character %q", domain.ErrValidation, r)
		}
	}

	return nil
}

// URLValidate пустой хттп/хттпс-адрес, который парсится.
func URLValidate(raw string) error {
	if strings.TrimSpace(raw) == "" {
		return fmt.Errorf("%w: url is empty", domain.ErrValidation)
	}

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return fmt.Errorf("%w: url is not valid", domain.ErrValidation)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: url scheme must be http or https", domain.ErrValidation)
	}

	if u.Host == "" {
		return fmt.Errorf("%w: url has no host", domain.ErrValidation)
	}

	return nil
}

// isExpected штатные ошибки, а не падение инфраструктуры.
func isExpected(err error) bool {
	return errors.Is(err, domain.ErrNotFound) ||
		errors.Is(err, domain.ErrConflict) ||
		errors.Is(err, domain.ErrValidation)
}
