package i18n

import (
	"errors"
)

var (
	ErrTerritoryNotFound = errors.New("i18n: territory not found")
)

type Territory struct {
	Code string
	NativeName string
	EnglishName string
}

func GetTerritory(code string) (*Territory, error) {
	if t, found := territories[code]; found {
		return t, nil
	}
	return nil, ErrTerritoryNotFound
}
