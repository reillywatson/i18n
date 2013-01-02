package i18n

import (
	"errors"
)

var (
	ErrTerritoryNotFound = errors.New("i18n: territory not found")

	allTerritories []*Territory
)

type Territory struct {
	Code string
	NativeName string
	EnglishName string
}

func init() {
	allTerritories = make([]*Territory, 0)
	for _, territory := range territories {
		allTerritories = append(allTerritories, territory)
	}
}

func GetTerritory(code string) (*Territory, error) {
	if t, found := territories[code]; found {
		return t, nil
	}
	return nil, ErrTerritoryNotFound
}

func Territories() []*Territory {
	return allTerritories
}
