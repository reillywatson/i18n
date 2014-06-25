package i18n

import (
	"reflect"
	"testing"
)

func TestAddressFormat(t *testing.T) {
	tests := []struct {
		Address  *Address
		Expected []string
	}{
		// #0: DE
		{
			Address: &Address{
				StreetAddress:   "Marienplatz 2a",
				ExtendedAddress: "",
				Locality:        "München",
				Region:          "Bayern",
				Country:         "DE",
				PostalCode:      "80331",
			},
			Expected: []string{
				"Marienplatz 2a",
				"80331 München",
				"Germany",
			},
		},
		// #1: DE
		{
			Address: &Address{
				StreetAddress:   "Marienplatz 2a",
				ExtendedAddress: "Rathausgasse",
				Locality:        "München",
				Region:          "Bayern",
				Country:         "DE",
				PostalCode:      "80331",
			},
			Expected: []string{
				"Marienplatz 2a",
				"Rathausgasse",
				"80331 München",
				"Germany",
			},
		},
		// #2: US
		{
			Address: &Address{
				StreetAddress:   "795 E Dragram Suite 5a",
				ExtendedAddress: "",
				Locality:        "Tucson",
				Region:          "AZ",
				Country:         "US",
				PostalCode:      "85705",
			},
			Expected: []string{
				"795 E Dragram Suite 5a",
				"Tucson AZ 85705",
				"United States",
			},
		},
		// #3: US
		{
			Address: &Address{
				StreetAddress:   "795 E Dragram Suite 5a",
				ExtendedAddress: "Building A",
				Locality:        "Tucson",
				Region:          "AZ",
				Country:         "US",
				PostalCode:      "85705",
			},
			Expected: []string{
				"795 E Dragram Suite 5a",
				"Building A",
				"Tucson AZ 85705",
				"United States",
			},
		},
		// #4: UK
		{
			Address: &Address{
				StreetAddress:   "49 Featherstone Street",
				ExtendedAddress: "",
				Locality:        "London",
				Region:          "",
				Country:         "GB",
				PostalCode:      "EC1Y 8SY",
			},
			Expected: []string{
				"49 Featherstone Street",
				"London",
				"EC1Y 8SY",
				"United Kingdom",
			},
		},
		// #5: UK
		{
			Address: &Address{
				StreetAddress:   "Oxford Road",
				ExtendedAddress: "Ardenham Court",
				Locality:        "Aylesbury",
				Region:          "Buckinghamshire",
				Country:         "GB",
				PostalCode:      "HP19 3EQ",
			},
			Expected: []string{
				"Oxford Road",
				"Ardenham Court",
				"Aylesbury",
				"Buckinghamshire",
				"HP19 3EQ",
				"United Kingdom",
			},
		},
		// #6: AT
		{
			Address: &Address{
				StreetAddress:   "Pazmaniteng 24-9",
				ExtendedAddress: "",
				Locality:        "Wien",
				Region:          "",
				Country:         "AT",
				PostalCode:      "1020",
			},
			Expected: []string{
				"Pazmaniteng 24-9",
				"A-1020 Wien",
				"Austria",
			},
		},
		// #7: CH
		{
			Address: &Address{
				StreetAddress:   "Kappelergasse 1",
				ExtendedAddress: "",
				Locality:        "Zürich",
				Region:          "",
				Country:         "CH",
				PostalCode:      "8022",
			},
			Expected: []string{
				"Kappelergasse 1",
				"8022 Zürich",
				"Switzerland",
			},
		},
		// #8: CH
		{
			Address: &Address{
				StreetAddress:   "Kappelergasse 1",
				ExtendedAddress: "",
				Locality:        "Zürich",
				Region:          "SZ",
				Country:         "CH",
				PostalCode:      "8022",
			},
			Expected: []string{
				"Kappelergasse 1",
				"8022 Zürich SZ",
				"Switzerland",
			},
		},
		// #9: ES
		{
			Address: &Address{
				StreetAddress:   "Calle Aduana, 29",
				ExtendedAddress: "",
				Locality:        "Madrid",
				Region:          "",
				Country:         "ES",
				PostalCode:      "28070",
			},
			Expected: []string{
				"Calle Aduana, 29",
				"28070 Madrid",
				"Spain",
			},
		},
		// #10: FR
		{
			Address: &Address{
				StreetAddress:   "27 Rue Pasteur",
				ExtendedAddress: "",
				Locality:        "Cabourg",
				Region:          "",
				Country:         "FR",
				PostalCode:      "14390",
			},
			Expected: []string{
				"27 Rue Pasteur",
				"14390 Cabourg",
				"France",
			},
		},
		// #11: AU
		{
			Address: &Address{
				StreetAddress:   "200 Broadway Av",
				ExtendedAddress: "",
				Locality:        "West Beach",
				Region:          "SA",
				Country:         "AU",
				PostalCode:      "5024",
			},
			Expected: []string{
				"200 Broadway Av",
				"West Beach SA 5024",
				"Australia",
			},
		},
		// #12: BR
		{
			Address: &Address{
				StreetAddress:   "Rua Visconde de Porto Seguro 1238",
				ExtendedAddress: "",
				Locality:        "São Paulo",
				Region:          "SP",
				Country:         "BR",
				PostalCode:      "04642-000",
			},
			Expected: []string{
				"Rua Visconde de Porto Seguro 1238",
				"São Paulo-SP",
				"04642-000",
				"Brazil",
			},
		},
		// #13: BR
		{
			Address: &Address{
				StreetAddress:   "Rua Visconde de Porto Seguro 1238",
				ExtendedAddress: "",
				Locality:        "São Paulo",
				Region:          "",
				Country:         "BR",
				PostalCode:      "04642-000",
			},
			Expected: []string{
				"Rua Visconde de Porto Seguro 1238",
				"São Paulo",
				"04642-000",
				"Brazil",
			},
		},
		// #14: CA
		{
			Address: &Address{
				StreetAddress:   "10-123 1/2 Main Street NW",
				ExtendedAddress: "",
				Locality:        "Montreal",
				Region:          "QC",
				Country:         "CA",
				PostalCode:      "H3Z 2Y7",
			},
			Expected: []string{
				"10-123 1/2 Main Street NW",
				"Montreal QC H3Z 2Y7",
				"Canada",
			},
		},
		// #15: MX
		{
			Address: &Address{
				StreetAddress:   "Col. Atlatilco",
				ExtendedAddress: "",
				Locality:        "Mexico",
				Region:          "D.F.",
				Country:         "MX",
				PostalCode:      "02860",
			},
			Expected: []string{
				"Col. Atlatilco",
				"02860 Mexico, D.F.",
				"Mexico",
			},
		},
		// #16: MX
		{
			Address: &Address{
				StreetAddress:   "8th Straco",
				ExtendedAddress: "",
				Locality:        "Puerto Vallarta",
				Region:          "",
				Country:         "MX",
				PostalCode:      "46800",
			},
			Expected: []string{
				"8th Straco",
				"46800 Puerto Vallarta",
				"Mexico",
			},
		},
	}

	for i, test := range tests {
		got := test.Address.FormattedParts()
		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf("test case #%d: got %q, expected %q", i, got, test.Expected)
		}
	}
}
