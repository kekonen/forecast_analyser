package fetchers

import (
	"forecast_analyser/lib/forecast"
)

type Fetcher interface {
	Fetch(cities []string, dailyForecasts chan forecast.DailyForecast, hourlyForecasts chan forecast.DailyForecast)
}
