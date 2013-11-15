package exchanges

import (
	"encoding/xml"
	"log"
	"net/http"
	"sync"
	"time"

	i18n "github.com/olivere/i18n-go"
)

var (
	ECBUrl = "http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

	ecbInstanceLock sync.Mutex
	ecbInstance     *ecb
)

func ECB() i18n.Exchanger {
	ecbInstanceLock.Lock()
	defer ecbInstanceLock.Unlock()

	// If we have, return the singleton
	if ecbInstance != nil {
		return ecbInstance
	}

	// Create a singleton
	ecbInstance = &ecb{
		rates: make(map[string]float64),
	}
	// Initial load
	ecbInstance.loadRates()
	// Start periodic updater
	go ecbInstance.updater()
	// Return instance
	return ecbInstance
}

type ecb struct {
	mu             sync.Mutex
	rates          map[string]float64
	lastDownloaded time.Time
}

func (ecb *ecb) GetRate(currency string) float64 {
	rate, found := ecb.rates[currency]
	if !found {
		return 0.0
	}
	return rate
}

func (ecb *ecb) updater() {
	ticker := time.Tick(1 * time.Hour)

	for {
		// Wait for next tick
		<-ticker

		ecb.loadRates()
	}
}

func (ecb *ecb) loadRates() error {
	ecb.mu.Lock()
	defer ecb.mu.Unlock()

	// Grab data from ECB via HTTP request
	resp, err := http.Get(ECBUrl)
	if err != nil {
		log.Printf("http.Get failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Decode XML
	var env envelope
	if err := xml.NewDecoder(resp.Body).Decode(&env); err != nil {
		log.Printf("decoding failed: %v", err)
		return err
	}

	// Update rates
	newRates := make(map[string]float64)
	for _, cube := range env.Cube.Rates {
		newRates[cube.Currency] = cube.Rate
	}
	ecb.rates = newRates
	ecb.lastDownloaded = time.Now()

	return nil
}

type envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Cube    *cube    `xml:"Cube>Cube"`
}

type cube struct {
	XMLName xml.Name `xml:"Cube"`
	Time    string   `xml:"time,attr"`
	Rates   []*rate  `xml:"Cube"`
}

type rate struct {
	XMLName  xml.Name `xml:"Cube"`
	Currency string   `xml:"currency,attr"`
	Rate     float64  `xml:"rate,attr"`
}
