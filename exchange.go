package i18n

import (
	_ "time"
)

var (
	Exchange Exchanger
)

type Exchanger interface {
	GetRate(string) float64
}

func init() {
	Exchange = new(nopExchanger)
}

type nopExchanger struct {
}

func (ex *nopExchanger) GetRate(currency string) float64 {
	return 1.0
}

// ExchangeTo converts a monetary value to another currency.
func (m *Money) ExchangeTo(currency string) *Money {
	rate := Exchange.GetRate(currency)
	y := m.Mulf(rate)
	y.C = currency
	return y
}
