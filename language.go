package i18n

import (
	"errors"
)

var (
	ErrLanguageNotFound = errors.New("i18n: language not found")
)

type Language struct {
	Code string
	NativeName string
	EnglishName string
}

func GetLanguage(code string) (*Language, error) {
	if l, found := languages[code]; found {
		return l, nil
	}
	return nil, ErrLanguageNotFound
}
