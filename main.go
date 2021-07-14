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
	db.AutoMigrate(&forecast.CurrentForecast{})
	db.AutoMigrate(&forecast.DailyForecast{})
	db.AutoMigrate(&forecast.HourlyForecast{})

	forecasts := make(chan interface{})
	cities := []string{
		"Berlin",
	}

	var wg sync.WaitGroup
	fetcher1 := fetchers.WttrFetcher{}
	fetcher2 := fetchers.NewOpenWeatherFetcher("f7d469c6ee94d8569b4f98bfe43fb4a1")

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	fetcher1.Fetch(cities, dailyForecasts, hourlyForecasts)
	// }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fetcher1.Fetch(cities, forecasts)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fetcher2.Fetch(cities, forecasts)
	}()

	currentForecasts := make([]forecast.CurrentForecast, 0)
	hourlyForecasts := make([]forecast.HourlyForecast, 0)
	dailyForecasts := make([]forecast.DailyForecast, 0)

	go func() {
		for f := range forecasts {
			switch v := f.(type) {
			case forecast.DailyForecast:
				fmt.Printf("dly: %v\n", v.Describe())
				dailyForecasts = append(dailyForecasts, v)
			case forecast.HourlyForecast:
				fmt.Printf("hly: %v\n", v.Describe())
				hourlyForecasts = append(hourlyForecasts, v)
			case forecast.CurrentForecast:
				fmt.Printf("Current: %v\n", v.Describe())
				currentForecasts = append(currentForecasts, v)
			default:
				fmt.Printf("I don't know about type %T!\n", v)
			}
		}
	}()

	wg.Wait()
	fmt.Println("Closing!")
	close(forecasts)

	db.CreateInBatches(currentForecasts, 100)
	db.CreateInBatches(hourlyForecasts, 100)
	db.CreateInBatches(dailyForecasts, 100)
}
