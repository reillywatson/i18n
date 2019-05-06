package i18n

// Currency represets all details about a currency.
type Currency struct {
	// Code is the 3-letter ISO code of the currency
	Code CurrencyCode
	// Symbol is the common symbol used for the currency, e.g. € for Euro.
	Symbol string
}

func CurrencyForCountryCode(countryCode string) CurrencyCode {
	for _, loc := range Locales {
		if loc.Territory == countryCode {
			return loc.CurrencyCode
		}
	}
	return ""
}
