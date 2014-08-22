package i18n

import (
	"testing"
)

func TestCurrencies(t *testing.T) {
	var tests = []struct {
		code         string
		found        bool
		expectedCode string
	}{
		/* 0 */ {"XYZ", false, ""},
		/* 1 */ {"EUR", true, "EUR"},
		/* 2 */ {"CHF", true, "CHF"},
		/* 3 */ {"USD", true, "USD"},
		/* 4 */ {"CNY", true, "CNY"},
	}

	for i, f := range tests {
		c, found := Currencies[f.code]
		if found != f.found {
			t.Fatalf("%d. expected currency code %v found flag to be %v, got %v", i, f.code, f.found, found)
		}
		if found {
			if c == nil {
				t.Fatalf("%d. expected currency to be != nil", i)
			}
			if f.expectedCode != c.Code {
				t.Errorf("%d. expected Code to be %v, got %v", i, f.expectedCode, c.Code)
			}
		}
	}
}
