package i18n

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestM(t *testing.T) {
	m := Money{}
	if m.M != 0 {
		t.Errorf("expected money amount to be %v, got %v", 0, m.M)
	}

	m = Money{-123, "EUR"}
	if m.M != -123 {
		t.Errorf("expected money amount to be %v, got %v", -123, m.M)
	}

	m = Money{123, "EUR"}
	if m.M != 123 {
		t.Errorf("expected money amount to be %v, got %v", 123, m.M)
	}
}

func TestGet(t *testing.T) {
	m := Money{123, "EUR"}
	if m.Get() != 1.23 {
		t.Errorf("expected money amount to be %v, got %v", 1.23, m.Get())
	}

	m = Money{12345, "EUR"}
	if m.Get() != 123.45 {
		t.Errorf("expected money amount to be %v, got %v", 123.45, m.Get())
	}
}

func TestAdd(t *testing.T) {
	m1 := Money{123, "EUR"}
	m2 := Money{678, "EUR"}
	m3 := m1.Add(m2)
	if m3.Get() != 8.01 {
		t.Errorf("expected money amount to be %v, got %v", 8.01, m3.Get())
	}
}

func TestAddEmpty(t *testing.T) {
	m1 := Money{}
	m2 := Money{123, "EUR"}
	m3 := m1.Add(m2)
	if m3.C != "EUR" {
		t.Errorf("expected currency to be EUR, got %v", m3.C)
	}
}

func TestAddYen(t *testing.T) {
	m1 := Money{123, JPY}
	m2 := Money{456, JPY}
	m3 := m1.Add(m2)
	if m3.Get() != 579 {
		t.Errorf("expected money amount to be 579, got %v", m3.Get())
	}
}

func TestMul(t *testing.T) {
	m1 := Money{123, "EUR"}
	m2 := Money{200, "EUR"}
	m3 := m1.Mul(m2)
	if m3.Get() != 2.46 {
		t.Errorf("expected money amount to be %v, got %v", 2.46, m3.Get())
	}
}

func TestMulYen(t *testing.T) {
	m1 := Money{123, JPY}
	m2 := Money{2, JPY}
	m3 := m1.Mul(m2)
	if m3.Get() != 246 {
		t.Errorf("expected money amount to be 246, got %v", m3.Get())
	}
}

func TestDiv(t *testing.T) {
	tests := []struct {
		name      string
		money1    Money
		money2    Money
		expResult Money
		expPanic  bool
	}{
		{
			name:      "Rational number with repeating decimals get truncated to 2 decimals",
			money1:    Money{1000, "EUR"},
			money2:    Money{300, "EUR"},
			expResult: Money{333, "EUR"},
		},
		{
			name:      "Rational number with repeating decimals get truncated and rounded to 2 decimals",
			money1:    Money{2000, "EUR"},
			money2:    Money{300, "EUR"},
			expResult: Money{667, "EUR"},
		},
		{
			name:      "Rational number with repeating decimals get truncated and rounded to 0 decimals",
			money1:    Money{2000, "JPY"},
			money2:    Money{300, "JPY"},
			expResult: Money{7, "JPY"},
		},
		{
			name:      "Rational number requiring truncation rounded correctly",
			money1:    Money{21252, "CAD"},
			money2:    MakeMoney("CAD", 24),
			expResult: Money{886, "CAD"},
		},
		{
			name:      "Rational number requiring truncation rounded correctly, second test",
			money1:    Money{22668, "USD"},
			money2:    MakeMoney("USD", 24),
			expResult: Money{945, "USD"},
		},
		{
			name:      "Rational number not requiring truncation is divided correctly",
			money1:    Money{90672, "CAD"},
			money2:    MakeMoney("CAD", 24),
			expResult: Money{3778, "CAD"},
		},
		// Testing for division by invalid numbers or amounts
		{
			name:      "Should panic when dividing by zero",
			money1:    Money{10, "CAD"},
			money2:    MakeMoney("CAD", 0),
			expResult: Money{0, "CAD"},
			expPanic:  true,
		},
	}

	for _, test := range tests {
		// Used for catching panics from dividing by 0!
		if test.expPanic {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s: Division by zero, should have paniced!", test.name)
				}
			}()
		}

		got := test.money1.Div(test.money2)
		if !test.expPanic && !reflect.DeepEqual(test.expResult, got) {
			t.Errorf("%s: expected money amount to be %v, got %v", test.name, test.expResult, got)
		}
	}
}

func TestMulf(t *testing.T) {
	var fixtures = []struct {
		money      Money
		multiplier float64
		expResult  float64
	}{
		{
			money:      Money{3390, "USD"},
			multiplier: .75,
			expResult:  25.43,
		},
		{
			money:      Money{5666, "USD"},
			multiplier: 1.0 / 6,
			expResult:  9.44,
		},
		{
			money:      Money{5667, "USD"},
			multiplier: 1.0 / 6,
			expResult:  9.45,
		},
		{
			money:      Money{5668, "USD"},
			multiplier: 1.0 / 6,
			expResult:  9.45,
		},
		{
			money:      Money{1234, "EUR"},
			multiplier: 2.26,
			expResult:  27.89,
		},
		{
			money:      Money{12345, "EUR"},
			multiplier: 2.26,
			expResult:  279,
		},
		{
			money:      Money{123456, "EUR"},
			multiplier: 2.26,
			expResult:  2790.11,
		},
		{
			money:      Money{123456, "EUR"},
			multiplier: 7 / 12.0,
			expResult:  720.16,
		},
		{
			money:      Money{-123456, "EUR"},
			multiplier: 7 / 12.0,
			expResult:  -720.16,
		},
	}

	for _, fixture := range fixtures {
		got := fixture.money.Mulf(fixture.multiplier).Get()
		if got != fixture.expResult {
			t.Errorf("expected money amount to be %v, got %v", fixture.expResult, got)
		}
	}
}

func TestSplit(t *testing.T) {
	var fixtures = []struct {
		money     Money
		chunks    int64
		expMoneys []Money
	}{
		{
			money:  Money{1, "CAD"},
			chunks: int64(1),
			expMoneys: []Money{
				{1, "CAD"},
			},
		},
		{
			money:  Money{2, "CAD"},
			chunks: int64(2),
			expMoneys: []Money{
				{1, "CAD"},
				{1, "CAD"},
			},
		},
		{
			money:  Money{700, "CAD"},
			chunks: int64(3),
			expMoneys: []Money{
				{234, "CAD"},
				{233, "CAD"},
				{233, "CAD"},
			},
		},
		{
			money:  Money{7, "CAD"},
			chunks: int64(3),
			expMoneys: []Money{
				{3, "CAD"},
				{2, "CAD"},
				{2, "CAD"},
			},
		},
		{
			money:  Money{99, "CAD"},
			chunks: int64(10),
			expMoneys: []Money{
				{10, "CAD"},
				{10, "CAD"},
				{10, "CAD"},
				{10, "CAD"},
				{10, "CAD"},
				{10, "CAD"},
				{10, "CAD"},
				{10, "CAD"},
				{10, "CAD"},
				{9, "CAD"},
			},
		},
	}

	for _, f := range fixtures {
		got := f.money.Split(f.chunks)
		if !reflect.DeepEqual(got, f.expMoneys) {
			t.Errorf("expected %s, got %s", f.expMoneys, got)
		}
	}
}

func TestMoneyStringer(t *testing.T) {
	var fixtures = []struct {
		m        Money
		expected string
	}{
		{Money{0, "EUR"}, "0.00 EUR"},
		{Money{1, "EUR"}, "0.01 EUR"},
		{Money{12, "EUR"}, "0.12 EUR"},
		{Money{123, "EUR"}, "1.23 EUR"},
		{Money{1234, "EUR"}, "12.34 EUR"},
		{Money{100000, "EUR"}, "1000.00 EUR"},
		{Money{123456, "EUR"}, "1234.56 EUR"},
		{Money{1234567, "EUR"}, "12345.67 EUR"},
		{Money{1234567890, "EUR"}, "12345678.90 EUR"},
		{Money{-1, "EUR"}, "-0.01 EUR"},
		{Money{-12, "EUR"}, "-0.12 EUR"},
		{Money{-123, "EUR"}, "-1.23 EUR"},
		{Money{-1234, "EUR"}, "-12.34 EUR"},
		{Money{-100000, "EUR"}, "-1000.00 EUR"},
		{Money{-123456, "EUR"}, "-1234.56 EUR"},
		{Money{-1234567, "EUR"}, "-12345.67 EUR"},
		{Money{-1234567890, "EUR"}, "-12345678.90 EUR"},
	}

	for _, f := range fixtures {
		got := fmt.Sprintf("%s", f.m)
		if got != f.expected {
			t.Errorf("expected %s, got %s", f.expected, got)
		}
	}
}

func TestMoneyFormat(t *testing.T) {
	var fixtures = []struct {
		m        Money
		locale   string
		expected string
	}{
		{Money{0, "EUR"}, "de", "0.00 EUR"},
		{Money{1, "EUR"}, "de", "0.01 EUR"},
		{Money{12, "EUR"}, "de", "0.12 EUR"},
		{Money{123, "EUR"}, "de", "1.23 EUR"},
		{Money{1234, "EUR"}, "de", "12.34 EUR"},
		{Money{123456, "EUR"}, "de", "1234.56 EUR"},
		{Money{1234567, "EUR"}, "de", "12345.67 EUR"},
		{Money{1234567890, "EUR"}, "de", "12345678.90 EUR"},
		{Money{-1, "EUR"}, "de", "-0.01 EUR"},
		{Money{-12, "EUR"}, "de", "-0.12 EUR"},
		{Money{-123, "EUR"}, "de", "-1.23 EUR"},
		{Money{-1234, "EUR"}, "de", "-12.34 EUR"},
		{Money{-123456, "EUR"}, "de", "-1234.56 EUR"},
		{Money{-1234567, "EUR"}, "de", "-12345.67 EUR"},
		{Money{-1234567890, "EUR"}, "de", "-12345678.90 EUR"},

		{Money{1, "EUR"}, "de", "0.01 EUR"},
		{Money{10, "EUR"}, "de", "0.10 EUR"},
		{Money{100, "EUR"}, "de", "1.00 EUR"},
		{Money{1000, "EUR"}, "de", "10.00 EUR"},
		{Money{10000, "EUR"}, "de", "100.00 EUR"},
		{Money{100000, "EUR"}, "de", "1000.00 EUR"},
		{Money{1000000, "EUR"}, "de", "10000.00 EUR"},
		{Money{10000000, "EUR"}, "de", "100000.00 EUR"},
		{Money{100000000, "EUR"}, "de", "1000000.00 EUR"},
		{Money{-1, "EUR"}, "de", "-0.01 EUR"},
		{Money{-10, "EUR"}, "de", "-0.10 EUR"},
		{Money{-100, "EUR"}, "de", "-1.00 EUR"},
		{Money{-1000, "EUR"}, "de", "-10.00 EUR"},
		{Money{-10000, "EUR"}, "de", "-100.00 EUR"},
		{Money{-100000, "EUR"}, "de", "-1000.00 EUR"},
		{Money{-1000000, "EUR"}, "de", "-10000.00 EUR"},
		{Money{-10000000, "EUR"}, "de", "-100000.00 EUR"},
		{Money{-100000000, "EUR"}, "de", "-1000000.00 EUR"},

		{Money{0, "EUR"}, "de_DE", "0,00 €"},
		{Money{1, "EUR"}, "de_DE", "0,01 €"},
		{Money{12, "EUR"}, "de_DE", "0,12 €"},
		{Money{123, "EUR"}, "de_DE", "1,23 €"},
		{Money{1234, "EUR"}, "de_DE", "12,34 €"},
		{Money{100000, "EUR"}, "de_DE", "1.000,00 €"},
		{Money{100100, "EUR"}, "de_DE", "1.001,00 €"},
		{Money{101000, "EUR"}, "de_DE", "1.010,00 €"},
		{Money{110000, "EUR"}, "de_DE", "1.100,00 €"},
		{Money{123456, "EUR"}, "de_DE", "1.234,56 €"},
		{Money{1234567, "EUR"}, "de_DE", "12.345,67 €"},
		{Money{1234567890, "EUR"}, "de_DE", "12.345.678,90 €"},
		{Money{-1, "EUR"}, "de_DE", "-0,01 €"},
		{Money{-12, "EUR"}, "de_DE", "-0,12 €"},
		{Money{-123, "EUR"}, "de_DE", "-1,23 €"},
		{Money{-1234, "EUR"}, "de_DE", "-12,34 €"},
		{Money{-100000, "EUR"}, "de_DE", "-1.000,00 €"},
		{Money{-123456, "EUR"}, "de_DE", "-1.234,56 €"},
		{Money{-1234567, "EUR"}, "de_DE", "-12.345,67 €"},
		{Money{-1234567890, "EUR"}, "de_DE", "-12.345.678,90 €"},

		{Money{1, "EUR"}, "de_DE", "0,01 €"},
		{Money{10, "EUR"}, "de_DE", "0,10 €"},
		{Money{100, "EUR"}, "de_DE", "1,00 €"},
		{Money{1000, "EUR"}, "de_DE", "10,00 €"},
		{Money{10000, "EUR"}, "de_DE", "100,00 €"},
		{Money{100000, "EUR"}, "de_DE", "1.000,00 €"},
		{Money{1000000, "EUR"}, "de_DE", "10.000,00 €"},
		{Money{10000000, "EUR"}, "de_DE", "100.000,00 €"},
		{Money{100000000, "EUR"}, "de_DE", "1.000.000,00 €"},
		{Money{-1, "EUR"}, "de_DE", "-0,01 €"},
		{Money{-10, "EUR"}, "de_DE", "-0,10 €"},
		{Money{-100, "EUR"}, "de_DE", "-1,00 €"},
		{Money{-1000, "EUR"}, "de_DE", "-10,00 €"},
		{Money{-10000, "EUR"}, "de_DE", "-100,00 €"},
		{Money{-100000, "EUR"}, "de_DE", "-1.000,00 €"},
		{Money{-1000000, "EUR"}, "de_DE", "-10.000,00 €"},
		{Money{-10000000, "EUR"}, "de_DE", "-100.000,00 €"},
		{Money{-100000000, "EUR"}, "de_DE", "-1.000.000,00 €"},

		{Money{0, "EUR"}, "de_AT", "€ 0,00"},
		{Money{1, "EUR"}, "de_AT", "€ 0,01"},
		{Money{12, "EUR"}, "de_AT", "€ 0,12"},
		{Money{123, "EUR"}, "de_AT", "€ 1,23"},
		{Money{1234, "EUR"}, "de_AT", "€ 12,34"},
		{Money{100000, "EUR"}, "de_AT", "€ 1.000,00"},
		{Money{123456, "EUR"}, "de_AT", "€ 1.234,56"},
		{Money{1234567, "EUR"}, "de_AT", "€ 12.345,67"},
		{Money{1234567890, "EUR"}, "de_AT", "€ 12.345.678,90"},
		{Money{-1, "EUR"}, "de_AT", "-€ 0,01"},
		{Money{-12, "EUR"}, "de_AT", "-€ 0,12"},
		{Money{-123, "EUR"}, "de_AT", "-€ 1,23"},
		{Money{-1234, "EUR"}, "de_AT", "-€ 12,34"},
		{Money{-100000, "EUR"}, "de_AT", "-€ 1.000,00"},
		{Money{-123456, "EUR"}, "de_AT", "-€ 1.234,56"},
		{Money{-1234567, "EUR"}, "de_AT", "-€ 12.345,67"},
		{Money{-1234567890, "EUR"}, "de_AT", "-€ 12.345.678,90"},

		{Money{1, "EUR"}, "de_AT", "€ 0,01"},
		{Money{10, "EUR"}, "de_AT", "€ 0,10"},
		{Money{100, "EUR"}, "de_AT", "€ 1,00"},
		{Money{1000, "EUR"}, "de_AT", "€ 10,00"},
		{Money{10000, "EUR"}, "de_AT", "€ 100,00"},
		{Money{100000, "EUR"}, "de_AT", "€ 1.000,00"},
		{Money{1000000, "EUR"}, "de_AT", "€ 10.000,00"},
		{Money{10000000, "EUR"}, "de_AT", "€ 100.000,00"},
		{Money{100000000, "EUR"}, "de_AT", "€ 1.000.000,00"},
		{Money{-1, "EUR"}, "de_AT", "-€ 0,01"},
		{Money{-10, "EUR"}, "de_AT", "-€ 0,10"},
		{Money{-100, "EUR"}, "de_AT", "-€ 1,00"},
		{Money{-1000, "EUR"}, "de_AT", "-€ 10,00"},
		{Money{-10000, "EUR"}, "de_AT", "-€ 100,00"},
		{Money{-100000, "EUR"}, "de_AT", "-€ 1.000,00"},
		{Money{-1000000, "EUR"}, "de_AT", "-€ 10.000,00"},
		{Money{-10000000, "EUR"}, "de_AT", "-€ 100.000,00"},
		{Money{-100000000, "EUR"}, "de_AT", "-€ 1.000.000,00"},

		{Money{0, "EUR"}, "de_CH", "€ 0.00"},
		{Money{1, "EUR"}, "de_CH", "€ 0.01"},
		{Money{12, "EUR"}, "de_CH", "€ 0.12"},
		{Money{123, "EUR"}, "de_CH", "€ 1.23"},
		{Money{1234, "EUR"}, "de_CH", "€ 12.34"},
		{Money{100000, "EUR"}, "de_CH", "€ 1'000.00"},
		{Money{123456, "EUR"}, "de_CH", "€ 1'234.56"},
		{Money{1234567, "EUR"}, "de_CH", "€ 12'345.67"},
		{Money{1234567890, "EUR"}, "de_CH", "€ 12'345'678.90"},
		{Money{-1, "EUR"}, "de_CH", "€-0.01"},
		{Money{-12, "EUR"}, "de_CH", "€-0.12"},
		{Money{-123, "EUR"}, "de_CH", "€-1.23"},
		{Money{-1234, "EUR"}, "de_CH", "€-12.34"},
		{Money{-100000, "EUR"}, "de_CH", "€-1'000.00"},
		{Money{-123456, "EUR"}, "de_CH", "€-1'234.56"},
		{Money{-1234567, "EUR"}, "de_CH", "€-12'345.67"},
		{Money{-1234567890, "EUR"}, "de_CH", "€-12'345'678.90"},

		{Money{1, "EUR"}, "de_CH", "€ 0.01"},
		{Money{10, "EUR"}, "de_CH", "€ 0.10"},
		{Money{100, "EUR"}, "de_CH", "€ 1.00"},
		{Money{1000, "EUR"}, "de_CH", "€ 10.00"},
		{Money{10000, "EUR"}, "de_CH", "€ 100.00"},
		{Money{100000, "EUR"}, "de_CH", "€ 1'000.00"},
		{Money{1000000, "EUR"}, "de_CH", "€ 10'000.00"},
		{Money{10000000, "EUR"}, "de_CH", "€ 100'000.00"},
		{Money{100000000, "EUR"}, "de_CH", "€ 1'000'000.00"},
		{Money{-1, "EUR"}, "de_CH", "€-0.01"},
		{Money{-10, "EUR"}, "de_CH", "€-0.10"},
		{Money{-100, "EUR"}, "de_CH", "€-1.00"},
		{Money{-1000, "EUR"}, "de_CH", "€-10.00"},
		{Money{-10000, "EUR"}, "de_CH", "€-100.00"},
		{Money{-100000, "EUR"}, "de_CH", "€-1'000.00"},
		{Money{-1000000, "EUR"}, "de_CH", "€-10'000.00"},
		{Money{-10000000, "EUR"}, "de_CH", "€-100'000.00"},
		{Money{-100000000, "EUR"}, "de_CH", "€-1'000'000.00"},

		{Money{0, "EUR"}, "en", "0.00 EUR"},
		{Money{1, "EUR"}, "en", "0.01 EUR"},
		{Money{12, "EUR"}, "en", "0.12 EUR"},
		{Money{123, "EUR"}, "en", "1.23 EUR"},
		{Money{1234, "EUR"}, "en", "12.34 EUR"},
		{Money{100000, "EUR"}, "en", "1000.00 EUR"},
		{Money{123456, "EUR"}, "en", "1234.56 EUR"},
		{Money{1234567, "EUR"}, "en", "12345.67 EUR"},
		{Money{1234567890, "EUR"}, "en", "12345678.90 EUR"},
		{Money{-1, "EUR"}, "en", "-0.01 EUR"},
		{Money{-12, "EUR"}, "en", "-0.12 EUR"},
		{Money{-123, "EUR"}, "en", "-1.23 EUR"},
		{Money{-1234, "EUR"}, "en", "-12.34 EUR"},
		{Money{-100000, "EUR"}, "en", "-1000.00 EUR"},
		{Money{-123456, "EUR"}, "en", "-1234.56 EUR"},
		{Money{-1234567, "EUR"}, "en", "-12345.67 EUR"},
		{Money{-1234567890, "EUR"}, "en", "-12345678.90 EUR"},

		{Money{0, "EUR"}, "en_US", "€0.00"},
		{Money{1, "EUR"}, "en_US", "€0.01"},
		{Money{12, "EUR"}, "en_US", "€0.12"},
		{Money{123, "EUR"}, "en_US", "€1.23"},
		{Money{1234, "EUR"}, "en_US", "€12.34"},
		{Money{100000, "EUR"}, "en_US", "€1,000.00"},
		{Money{123456, "EUR"}, "en_US", "€1,234.56"},
		{Money{1234567, "EUR"}, "en_US", "€12,345.67"},
		{Money{1234567890, "EUR"}, "en_US", "€12,345,678.90"},
		{Money{-1, "EUR"}, "en_US", "(€0.01)"},
		{Money{-12, "EUR"}, "en_US", "(€0.12)"},
		{Money{-123, "EUR"}, "en_US", "(€1.23)"},
		{Money{-1234, "EUR"}, "en_US", "(€12.34)"},
		{Money{-100000, "EUR"}, "en_US", "(€1,000.00)"},
		{Money{-123456, "EUR"}, "en_US", "(€1,234.56)"},
		{Money{-1234567, "EUR"}, "en_US", "(€12,345.67)"},
		{Money{-1234567890, "EUR"}, "en_US", "(€12,345,678.90)"},

		{Money{0, "EUR"}, "fr", "0.00 EUR"},
		{Money{1, "EUR"}, "fr", "0.01 EUR"},
		{Money{12, "EUR"}, "fr", "0.12 EUR"},
		{Money{123, "EUR"}, "fr", "1.23 EUR"},
		{Money{1234, "EUR"}, "fr", "12.34 EUR"},
		{Money{100000, "EUR"}, "fr", "1000.00 EUR"},
		{Money{123456, "EUR"}, "fr", "1234.56 EUR"},
		{Money{1234567, "EUR"}, "fr", "12345.67 EUR"},
		{Money{1234567890, "EUR"}, "fr", "12345678.90 EUR"},
		{Money{-1, "EUR"}, "fr", "-0.01 EUR"},
		{Money{-12, "EUR"}, "fr", "-0.12 EUR"},
		{Money{-123, "EUR"}, "fr", "-1.23 EUR"},
		{Money{-1234, "EUR"}, "fr", "-12.34 EUR"},
		{Money{-100000, "EUR"}, "fr", "-1000.00 EUR"},
		{Money{-123456, "EUR"}, "fr", "-1234.56 EUR"},
		{Money{-1234567, "EUR"}, "fr", "-12345.67 EUR"},
		{Money{-1234567890, "EUR"}, "fr", "-12345678.90 EUR"},

		{Money{0, "EUR"}, "zh", "0.00 EUR"},
		{Money{1, "EUR"}, "zh", "0.01 EUR"},
		{Money{12, "EUR"}, "zh", "0.12 EUR"},
		{Money{123, "EUR"}, "zh", "1.23 EUR"},
		{Money{1234, "EUR"}, "zh", "12.34 EUR"},
		{Money{123456, "EUR"}, "zh", "1234.56 EUR"},
		{Money{1234567, "EUR"}, "zh", "12345.67 EUR"},
		{Money{1234567890, "EUR"}, "zh", "12345678.90 EUR"},
		{Money{-1, "EUR"}, "zh", "-0.01 EUR"},
		{Money{-12, "EUR"}, "zh", "-0.12 EUR"},
		{Money{-123, "EUR"}, "zh", "-1.23 EUR"},
		{Money{-1234, "EUR"}, "zh", "-12.34 EUR"},
		{Money{-123456, "EUR"}, "zh", "-1234.56 EUR"},
		{Money{-1234567, "EUR"}, "zh", "-12345.67 EUR"},
		{Money{-1234567890, "EUR"}, "zh", "-12345678.90 EUR"},

		{Money{0, "CNY"}, "zh_CN", "¥0.00"},
		{Money{1, "CNY"}, "zh_CN", "¥0.01"},
		{Money{12, "CNY"}, "zh_CN", "¥0.12"},
		{Money{123, "CNY"}, "zh_CN", "¥1.23"},
		{Money{1234, "CNY"}, "zh_CN", "¥12.34"},
		{Money{100000, "CNY"}, "zh_CN", "¥1,000.00"},
		{Money{123456, "CNY"}, "zh_CN", "¥1,234.56"},
		{Money{1234567, "CNY"}, "zh_CN", "¥12,345.67"},
		{Money{1234567890, "CNY"}, "zh_CN", "¥12,345,678.90"},
		{Money{-1, "CNY"}, "zh_CN", "¥-0.01"},
		{Money{-12, "CNY"}, "zh_CN", "¥-0.12"},
		{Money{-123, "CNY"}, "zh_CN", "¥-1.23"},
		{Money{-1234, "CNY"}, "zh_CN", "¥-12.34"},
		{Money{-100000, "CNY"}, "zh_CN", "¥-1,000.00"},
		{Money{-123456, "CNY"}, "zh_CN", "¥-1,234.56"},
		{Money{-1234567, "CNY"}, "zh_CN", "¥-12,345.67"},
		{Money{-1234567890, "CNY"}, "zh_CN", "¥-12,345,678.90"},

		{Money{1234567890, "USD"}, "en_US", "$12,345,678.90"},
		{Money{-1234567890, "USD"}, "en_US", "($12,345,678.90)"},
		{Money{1234567890, "USD"}, "de_DE", "12.345.678,90 $"},
		{Money{-1234567890, "USD"}, "de_DE", "-12.345.678,90 $"},
		{Money{1234567890, "USD"}, "de_CH", "$ 12'345'678.90"},
		{Money{-1234567890, "USD"}, "de_CH", "$-12'345'678.90"},
		{Money{1234567890, "USD"}, "zh_CN", "$12,345,678.90"},
		{Money{-1234567890, "USD"}, "zh_CN", "$-12,345,678.90"},
		{Money{1234567890, "GBP"}, "en_GB", "£12,345,678.90"},
		{Money{-1234567890, "GBP"}, "en_GB", "-£12,345,678.90"},
		{Money{1234567890, "GBP"}, "de_DE", "12.345.678,90 £"},
		{Money{-1234567890, "GBP"}, "de_DE", "-12.345.678,90 £"},
		{Money{1234567890, "GBP"}, "de_CH", "£ 12'345'678.90"},
		{Money{-1234567890, "GBP"}, "de_CH", "£-12'345'678.90"},
		{Money{1234567890, "GBP"}, "zh_CN", "£12,345,678.90"},
		{Money{-1234567890, "GBP"}, "zh_CN", "£-12,345,678.90"},
		{Money{1234567890, "HUF"}, "hu_HU", "12 345 678,90 Ft"},
		{Money{-1234567890, "HUF"}, "hu_HU", "-12 345 678,90 Ft"},
		{Money{1234567890, "HUF"}, "de_DE", "12.345.678,90 Ft"},
		{Money{-1234567890, "HUF"}, "de_DE", "-12.345.678,90 Ft"},
		{Money{1234567890, "HUF"}, "de_CH", "Ft 12'345'678.90"},
		{Money{-1234567890, "HUF"}, "de_CH", "Ft-12'345'678.90"},
		{Money{1234567890, "HUF"}, "zh_CN", "Ft12,345,678.90"},
		{Money{-1234567890, "HUF"}, "zh_CN", "Ft-12,345,678.90"},
		{Money{1234567890, "JPY"}, "ja_JP", "¥1,234,567,890"},
		{Money{-1234567890, "JPY"}, "ja_JP", "-¥1,234,567,890"},
		{Money{1234567890, "JPY"}, "de_DE", "12.345.678,90 ¥"},
		{Money{-1234567890, "JPY"}, "de_DE", "-12.345.678,90 ¥"},
		{Money{1234567890, "JPY"}, "de_CH", "¥ 12'345'678.90"},
		{Money{-1234567890, "JPY"}, "de_CH", "¥-12'345'678.90"},
		{Money{1234567890, "SEK"}, "se_SE", "12.345.678,90 kr"},
		{Money{-1234567890, "SEK"}, "se_SE", "-12.345.678,90 kr"},
		{Money{1234567890, "SEK"}, "de_DE", "12.345.678,90 kr"},
		{Money{-1234567890, "SEK"}, "de_DE", "-12.345.678,90 kr"},
		{Money{1234567890, "SEK"}, "de_CH", "kr 12'345'678.90"},
		{Money{-1234567890, "SEK"}, "de_CH", "kr-12'345'678.90"},
		{Money{1234567890, "SEK"}, "zh_CN", "kr12,345,678.90"},
		{Money{-1234567890, "SEK"}, "zh_CN", "kr-12,345,678.90"},
	}

	for _, f := range fixtures {
		got := f.m.Format(f.locale)
		if got != f.expected {
			t.Errorf("expected %s, got %s (locale: %s)", f.expected, got, f.locale)
		}
	}
}

func TestJSONUnmarshal(t *testing.T) {
	tests := []struct {
		JSON    string
		Money   Money
		Success bool
	}{
		{`{"C": "CAD","M": 500}`, Money{C: "CAD", M: 500}, true},
		{`{"C": "CAD","F": 5}`, Money{C: "CAD", M: 500}, true},
		{`{"C": "CAD","F": 5, "M": 1}`, Money{C: "CAD", M: 500}, true},
	}
	for _, test := range tests {
		var m Money
		err := json.Unmarshal([]byte(test.JSON), &m)
		if (err == nil) != test.Success {
			t.Errorf("Expected success %t, got error %v", test.Success, err)
			continue
		}
		if !reflect.DeepEqual(test.Money, m) {
			t.Errorf("Expected %+v, got %+v. JSON: %s", test.Money, m, test.JSON)
		}
	}
}

func TestUnitAsString(t *testing.T) {
	var fixtures = []struct {
		m        Money
		expected string
	}{
		{Money{0, "EUR"}, "0.00"},
		{Money{1, "EUR"}, "0.01"},
		{Money{12, "EUR"}, "0.12"},
		{Money{123, "EUR"}, "1.23"},
		{Money{1234, "EUR"}, "12.34"},
		{Money{100000, "EUR"}, "1000.00"},
		{Money{123456, "EUR"}, "1234.56"},
		{Money{1234567, "EUR"}, "12345.67"},
		{Money{1234567890, "EUR"}, "12345678.90"},
		{Money{-1, "EUR"}, "-0.01"},
		{Money{-12, "EUR"}, "-0.12"},
		{Money{-123, "EUR"}, "-1.23"},
		{Money{-1234, "EUR"}, "-12.34"},
		{Money{-100000, "EUR"}, "-1000.00"},
		{Money{-123456, "EUR"}, "-1234.56"},
		{Money{-1234567, "EUR"}, "-12345.67"},
		{Money{-1234567890, "EUR"}, "-12345678.90"},
	}

	for _, f := range fixtures {
		got := f.m.UnitAsString()
		if got != f.expected {
			t.Errorf("expected %s, got %s", f.expected, got)
		}
	}
}
