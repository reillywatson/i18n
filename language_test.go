package i18n

import (
	"testing"
)

func TestLanguages(t *testing.T) {
	var tests = []struct {
		code                string
		err                 error
		expectedEnglishName string
		expectedNativeName  string
	}{
		{"xy", ErrLanguageNotFound, "", ""},
		{"de", nil, "German", "Deutsch"},
		{"en", nil, "English", "English"},
		{"fr", nil, "French", "français"},
		{"es", nil, "Spanish", "español"},
	}

	for _, f := range tests {
		l, err := GetLanguage(f.code)
		if err != f.err {
			t.Fatalf("expected error to be %v, got %v", f.err, err)
		}
		if f.err == nil {
			if l == nil {
				t.Fatalf("expected language to be != nil")
			}
			if f.expectedEnglishName != l.EnglishName {
				t.Errorf("expected EnglishName to be %v, got %v", f.expectedEnglishName, l.EnglishName)
			}
			if f.expectedNativeName != l.NativeName {
				t.Errorf("expected NativeName to be %v, got %v", f.expectedNativeName, l.NativeName)
			}
		}
	}
}

func TestAllLanguages(t *testing.T) {
	all := Languages()
	if len(all) <= 0 {
		t.Errorf("expected list of languages, got none")
	}
}
