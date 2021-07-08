package main

import (
	"fmt"
	"time"
)

type Forecast struct {
	CollectedAt time.Time
	TargetAt    time.Time
	Type        int32 // "now", "hourly", "daily"
	Temperature float32
	Condition   string
	Location    string
	Source      string
}

type Fetcher interface {
	fetch(cities []string, forecasts chan Forecast)
}

type WttrFetcher struct {
}

func (f *WttrFetcher) fetch(cities []string, forecasts chan Forecast) {
	for _, city := range cities {
		if f.hasCity(&city) {
			f.fetchCity(city, forecasts)
		}
	}
}

func (f *WttrFetcher) hasCity(city *string) bool {
	return *city == "Berlin"
}

func (f *WttrFetcher) fetchCity(city string, forecasts chan Forecast) {
	forecasts <- Forecast{}
}

func main() {
	forecasts := make(chan Forecast)
	cities := []string{
		"Berlin",
	}
	fetcher1 := WttrFetcher{}

	go func() {
		fetcher1.fetch(cities, forecasts)
		close(forecasts)
	}()

	for forecast := range forecasts {
		fmt.Printf("%v\n", forecast)
	}
}
