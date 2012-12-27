package i18n

import (
	"testing"
)

func TestCurrencies(t *testing.T) {
	var tests = []struct {
		code           string
		err            error
		expectedCode   string
	}{
		{"XYZ", ErrCurrencyNotFound, ""},
		{"EUR", nil, "EUR"},
		{"CHF", nil, "CHF"},
		{"USD", nil, "USD"},
	}

	for _, f := range tests {
		c, err := GetCurrency(f.code)
		if err != f.err {
			t.Fatalf("expected error to be %v, got %v", f.err, err)
		}
		if err == nil {
			if c == nil {
				t.Fatalf("expected currency to be != nil")
			}
			if f.expectedCode != c.Code {
				t.Errorf("expected Code to be %v, got %v", f.expectedCode, c.Code)
			}
		}
	}
}
