package fetchers

import (
	"fmt"
	"forecast_analyser/lib/forecast"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

// https://openweathermap.org/api/one-call-api

var OWApiKey = "f7d469c6ee94d8569b4f98bfe43fb4a1"

type OpenWeatherFetcher struct {
}

func (f *OpenWeatherFetcher) Fetch(cities []string, dailyForecasts chan forecast.DailyForecast, hourlyForecasts chan forecast.HourlyForecast) {
	for _, city := range cities {
		if f.hasCity(&city) {
			f.fetchCity(city, dailyForecasts, hourlyForecasts)
		}
	}
}

func (f *OpenWeatherFetcher) source() string {
	return "Open Weather"
}

func (f *OpenWeatherFetcher) hasCity(city *string) bool {
	return *city == "Berlin"
}

func (f *OpenWeatherFetcher) getCityUrl(city *string) string {
	return fmt.Sprintf("http://wttr.in/%s?format=j1", *city)
}

func (f *OpenWeatherFetcher) fetchCity(city string, dailyForecasts chan forecast.DailyForecast, hourlyForecasts chan forecast.HourlyForecast) {
	url := f.getCityUrl(&city)
	// fmt.Printf("Wttr.in fetching: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Printf("Error happened: %v\n", err)
	}
	defer resp.Body.Close()
	jsonStr, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Printf("Error happened: %v\n", err)
	}

	localObsDateTime := gjson.Get(string(jsonStr), "current_condition.0.localObsDateTime").String()
	location, err := time.LoadLocation("Europe/Berlin")
	locatlObesrvationTime, _ := time.ParseInLocation("2006-01-02 3:04 PM", localObsDateTime, location)

	tempCStr := gjson.Get(string(jsonStr), "current_condition.0.temp_C").String()
	tempC, _ := strconv.Atoi(tempCStr)
	weatherCode := gjson.Get(string(jsonStr), "current_condition.0.weatherCode").Int()

	currentCondition := forecast.HourlyForecast{
		TargetAt:    locatlObesrvationTime,
		Temperature: float32(tempC),
		Condition:   parseWeatherCode(weatherCode),
		Location:    city,
		Source:      f.source(),
	}

	hourlyForecasts <- currentCondition
	weatherDays := gjson.Get(string(jsonStr), "weather")
	weatherDays.ForEach(func(key gjson.Result, value gjson.Result) bool {
		dateStr := value.Get("date")
		// fmt.Printf("Date: %v\n", dateStr)

		tempMin := value.Get("mintempC").Int()
		tempMax := value.Get("maxtempC").Int()

		date, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", dateStr), location)

		dailyForecast := forecast.DailyForecast{
			TargetAt:       date,
			TemperatureMin: float32(tempMin),
			TemperatureMax: float32(tempMax),
			Condition:      parseWeatherCode(weatherCode),
			Location:       city,
			Source:         f.source(),
		}

		dailyForecasts <- dailyForecast

		value.Get("hourly").ForEach(func(key gjson.Result, value gjson.Result) bool {
			timeStr := value.Get("time").String()
			if len(timeStr) < 2 {
				timeStr = "0000"
			} else if len(timeStr) < 4 {
				timeStr = fmt.Sprintf("0%s", timeStr)
			}
			tempC := value.Get("tempC").Float()
			weatherCode := value.Get("weatherCode").Int()

			// fmt.Printf(" * %v: %v %v\n", timeStr, tempC, parseWeatherCode(weatherCode))

			fullDate, err := time.ParseInLocation("2006-01-02 1504", fmt.Sprintf("%s %s", dateStr, timeStr), location)
			if err != nil {
				fmt.Println(err)
			}

			hourlyForecast := forecast.HourlyForecast{
				TargetAt:    fullDate,
				Temperature: float32(tempC),
				Condition:   parseWeatherCode(weatherCode),
				Location:    city,
				Source:      f.source(),
			}

			// fmt.Printf("%v\n", hourlyForecast)
			hourlyForecasts <- hourlyForecast

			return true
		})
		return true
	})
}
