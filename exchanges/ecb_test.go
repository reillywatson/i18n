package exchanges

import (
	"testing"

	i18n "github.com/olivere/i18n-go"
)

func TestECBRatesLoader(t *testing.T) {
	exchange := ECB()
	if ecbInstance == nil {
		t.Fatal("expected ECB singleton to be initialized")
		t.FailNow()
	}
	ecbInstance.mu.Lock()
	ecbInstance.mu.Unlock()
	if ecbInstance.lastDownloaded.IsZero() {
		t.Fatal("expected rates to be downloaded")
	}

	if len(ecbInstance.rates) == 0 {
		t.Error("expected to have some exchange rates loaded")
	}
	_, found := ecbInstance.rates["USD"]
	if !found {
		t.Error("expected to find exchange rate for USD")
	}
	usd := exchange.GetRate("USD")
	if usd <= 0.0 {
		t.Errorf("expected EUR -> USD rate of greater than 0, got: %v", usd)
	}
}

func TestECBMoneyExchangeTo(t *testing.T) {
	i18n.Exchange = ECB()
	in := &i18n.Money{100, "EUR"}
	out := in.Clone().ExchangeTo("USD")
	if out.M <= in.M {
		t.Errorf("expected a higher amount for USD than for EUR, "+
			"but %s is exchanged to %s", in, out)
	}
	if out.C != "USD" {
		t.Errorf("expected currency USD; got: %s", out.C)
	}
}
