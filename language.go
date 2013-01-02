package i18n

import (
	"errors"
)

var (
	ErrLanguageNotFound = errors.New("i18n: language not found")

	allLanguages []*Language
)

type Language struct {
	Code string
	NativeName string
	EnglishName string
}

func init() {
	allLanguages = make([]*Language, 0)
	for _, lang := range languages {
		allLanguages = append(allLanguages, lang)
	}
}

func GetLanguage(code string) (*Language, error) {
	if l, found := languages[code]; found {
		return l, nil
	}
	return nil, ErrLanguageNotFound
}

func Languages() []*Language {
	return allLanguages
}
