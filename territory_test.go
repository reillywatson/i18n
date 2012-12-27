package i18n

import (
	"testing"
)

func TestTerritories(t *testing.T) {
	var tests = []struct {
		code                string
		err                 error
		expectedEnglishName string
		expectedNativeName  string
	}{
		{"XY", ErrTerritoryNotFound, "", ""},
		{"DE", nil, "Germany", "Deutschland"},
		// BUG(oe): Wrong native name, at least from a German standpoint ;-)
		{"CH", nil, "Switzerland", "Svizra"},
		// BUG(oe): Wrong native name, at least from a German standpoint ;-)
		{"GB", nil, "United Kingdom", "y Deyrnas Unedig"},
		{"US", nil, "United States", "United States"},
	}

	for _, f := range tests {
		c, err := GetTerritory(f.code)
		if err != f.err {
			t.Fatalf("expected error to be %v, got %v", f.err, err)
		}
		if f.err == nil {
			if c == nil {
				t.Fatalf("expected country to be != nil")
			}
			if f.expectedEnglishName != c.EnglishName {
				t.Errorf("expected EnglishName to be %v, got %v", f.expectedEnglishName, c.EnglishName)
			}
			if f.expectedNativeName != c.NativeName {
				t.Errorf("expected NativeName to be %v, got %v", f.expectedNativeName, c.NativeName)
			}
		}
	}
}
