package main

import (
	"fmt"
	"sync"

	"forecast_analyser/lib/fetchers"
	"forecast_analyser/lib/forecast"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// type Forecast struct {
// 	gorm.Model
// 	TargetAt    time.Time
// 	Type        int32 // "now", "hourly", "daily"
// 	Temperature float32
// 	Condition   string
// 	Location    string
// 	Source      string
// }

// type Fetcher interface {
// 	fetch(cities []string, forecasts chan Forecast)
// }

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&forecast.DailyForecast{})
	db.AutoMigrate(&forecast.HourlyForecast{})

	dailyForecasts := make(chan forecast.DailyForecast)
	hourlyForecasts := make(chan forecast.HourlyForecast)
	cities := []string{
		"Berlin",
	}

	var wg sync.WaitGroup
	// fetcher1 := fetchers.WttrFetcher{}
	fetcher2 := fetchers.WttrFetcher{}

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	fetcher1.Fetch(cities, dailyForecasts, hourlyForecasts)
	// }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fetcher2.Fetch(cities, dailyForecasts, hourlyForecasts)
	}()

	go func() {
		for forecast := range dailyForecasts {
			fmt.Printf("%v\n", forecast.Describe())
		}
	}()

	go func() {
		for forecast := range hourlyForecasts {
			fmt.Printf("%v\n", forecast.Describe())
		}
	}()

	wg.Wait()
	fmt.Println("Closing!")
	close(dailyForecasts)
}
