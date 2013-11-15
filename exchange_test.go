package i18n

import (
	"testing"
)

type fixedExchanger struct {
	Rate float64
}

func (e *fixedExchanger) GetRate(currency string) float64 {
	return e.Rate
}

func TestExchangeTo(t *testing.T) {
	tests := []struct {
		Exchange Exchanger
		Input    *Money
		Currency string
		Output   *Money
	}{
		{
			Exchange: &fixedExchanger{1.0},
			Input:    &Money{0, "EUR"},
			Currency: "USD",
			Output:   &Money{0, "USD"},
		},
		{
			Exchange: &fixedExchanger{1.0},
			Input:    &Money{1, "EUR"},
			Currency: "USD",
			Output:   &Money{1, "USD"},
		},
		{
			Exchange: &fixedExchanger{2.0},
			Input:    &Money{1, "EUR"},
			Currency: "USD",
			Output:   &Money{2, "USD"},
		},
		{
			Exchange: &fixedExchanger{2.5},
			Input:    &Money{10, "EUR"},
			Currency: "USD",
			Output:   &Money{25, "USD"},
		},
		{
			Exchange: &fixedExchanger{2.5},
			Input:    &Money{10, "EUR"},
			Currency: "USD",
			Output:   &Money{25, "USD"},
		},
	}

	for _, test := range tests {
		Exchange = test.Exchange
		got := test.Input.Clone().ExchangeTo(test.Currency)
		if got.M != test.Output.M || got.C != test.Output.C {
			t.Errorf("expected %v; got: %v", test.Output, got)
		}
	}
}
