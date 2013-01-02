package i18n

import (
	"testing"
)

func TestLocales(t *testing.T) {
	var tests = []struct {
		code              string
		err               error
		expectedLanguage  string
		expectedLangISO3  string
		expectedTerritory string
		expectedCurrCode  string
	}{
		{"xy",    ErrLocaleNotFound, "", "", "", ""},
		{"de",    ErrLocaleNotFound, "", "", "", ""},
		{"en",    ErrLocaleNotFound, "", "", "", ""},
		{"fr",    ErrLocaleNotFound, "", "", "", ""},
		{"de_AT", nil, "de", "deu", "AT", "EUR"},
		{"de_CH", nil, "de", "deu", "CH", "CHF"},
		{"de_DE", nil, "de", "deu", "DE", "EUR"},
		{"en_CA", nil, "en", "eng", "CA", "CAD"},
		{"en_GB", nil, "en", "eng", "GB", "GBP"},
		{"en_US", nil, "en", "eng", "US", "USD"},
		{"fr_CH", nil, "fr", "fra", "CH", "CHF"},
		{"fr_FR", nil, "fr", "fra", "FR", "EUR"},
	}

	for _, f := range tests {
		loc, err := GetLocale(f.code)
		if err != f.err {
			t.Fatalf("expected error to be %v, got %v", f.err, err)
		}
		if err == nil {
			if loc == nil {
				t.Fatalf("expected locale to be != nil")
			}
			if f.expectedLanguage != loc.Language {
				t.Errorf("expected Language to be %v, got %v", f.expectedLanguage, loc.Language)
			}
			if f.expectedLangISO3 != loc.LanguageISO3 {
				t.Errorf("expected LanguageISO3 to be %v, got %v", f.expectedLangISO3, loc.LanguageISO3)
			}
			if f.expectedTerritory != loc.Territory {
				t.Errorf("expected Territory to be %v, got %v", f.expectedTerritory, loc.Territory)
			}
			if f.expectedCurrCode != loc.CurrencyCode {
				t.Errorf("expected CurrencyCode to be %v, got %v", f.expectedCurrCode, loc.CurrencyCode)
			}
		}
	}
}

func TestAllLocales(t *testing.T) {
	all := Locales()
	if len(all) <= 0 {
		t.Errorf("expected list of locales, got none")
	}
}
