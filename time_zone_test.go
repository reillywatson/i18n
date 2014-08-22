package i18n

import (
	"testing"
	"time"
)

func TestTimeZones(t *testing.T) {
	var tests = []struct {
		code  string
		found bool
	}{
		/* 0 */ {"Europe/Berlin", true},
		/* 1 */ {"Mars/Rover", false},
	}

	for i, test := range tests {
		var found bool
		for _, name := range TimeZones {
			if name == test.code {
				found = true
				break
			}
		}
		if found != test.found {
			t.Fatalf("%d. expected time zone %s found flag to be %v, got %v", i, test.code, test.found, found)
		}
	}
}

func TestTimeZoneAvailableInGo(t *testing.T) {
	for _, tz := range TimeZones {
		location, err := time.LoadLocation(tz)
		if err != nil {
			t.Fatalf("expected to load time zone %q, got %v", tz, err)
		}
		if location == nil {
			t.Errorf("expected location to not be nil with time zone %q", tz)
		}
	}
}
