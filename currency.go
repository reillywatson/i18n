package i18n

import (
	"errors"
)

var (
	ErrCurrencyNotFound = errors.New("i18n: currency not found")
)

type Currency struct {
	Code   string
	Symbol string
}

func GetCurrency(currency string) (*Currency, error) {
	if c, found := currencies[currency]; found {
		return c, nil
	}
	return nil, ErrCurrencyNotFound
}
