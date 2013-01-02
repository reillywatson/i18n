package i18n

import (
	"errors"
)

var (
	ErrCurrencyNotFound = errors.New("i18n: currency not found")

	allCurrencies []*Currency
)

type Currency struct {
	Code   string
	Symbol string
}

func init() {
	allCurrencies = make([]*Currency, 0)
	for _, currency := range currencies {
		allCurrencies = append(allCurrencies, currency)
	}
}

func GetCurrency(currency string) (*Currency, error) {
	if c, found := currencies[currency]; found {
		return c, nil
	}
	return nil, ErrCurrencyNotFound
}

func Currencies() []*Currency {
	return allCurrencies
}
