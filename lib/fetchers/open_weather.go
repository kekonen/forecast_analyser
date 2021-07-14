package fetchers

import (
	"fmt"
	"forecast_analyser/lib/forecast"
	"io"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

// https://openweathermap.org/api/one-call-api

var OWApiKey = "f7d469c6ee94d8569b4f98bfe43fb4a1"

type OpenWeatherFetcher struct {
	key string
}

func NewOpenWeatherFetcher(key string) OpenWeatherFetcher {
	return OpenWeatherFetcher{
		key,
	}
}

func (f *OpenWeatherFetcher) Fetch(cities []string, forecasts chan interface{}) {
	for _, city := range cities {
		if f.hasCity(&city) {
			f.fetchCity(city, forecasts)
		}
	}
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

func (f *OpenWeatherFetcher) source() string {
	return "Open Weather"
}

func (f *OpenWeatherFetcher) hasCity(city *string) bool {
	return *city == "Berlin"
}
func city2latlon(city *string) *Coordinates {
	if *city == "Berlin" {
		return &Coordinates{Latitude: 52.52, Longitude: 13.4}
	} else {
		return nil
	}
}

func (f *OpenWeatherFetcher) getCityUrl(city *string) string {
	coordinates := city2latlon(city)
	if coordinates != nil {
		return fmt.Sprintf(
			"https://api.openweathermap.org/data/2.5/onecall?lat=%f&lon=%f&units=metric&exclude=minutely&appid=%s",
			coordinates.Latitude,
			coordinates.Longitude,
			f.key,
		)
	} else {
		return ""
	}
}

func (f *OpenWeatherFetcher) fetchCity(city string, forecasts chan interface{}) {
	url := f.getCityUrl(&city)
	// fmt.Printf("Wttr.in fetching: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		// handle error
		fmt.Printf("Error happened: %v\n", err)
	}
	defer resp.Body.Close()
	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Printf("Error happened: %v\n", err)
	}

	jsonStr := string(jsonBytes)

	// localObsDateTime := gjson.Get(jsonStr, "current_condition.0.localObsDateTime").String()
	location, err := time.LoadLocation("Europe/Berlin")
	// locatlObesrvationTime, _ := time.ParseInLocation("2006-01-02 3:04 PM", localObsDateTime, location)

	temp := gjson.Get(jsonStr, "current.temp").Float()
	// tempC, _ := strconv.Atoi(tempCStr)
	weatherId := gjson.Get(jsonStr, "current.weather.id").Int()

	currentCondition := forecast.CurrentForecast{
		Temperature: float32(temp),
		Condition:   f.parseWeatherId(weatherId),
		Location:    city,
		Source:      f.source(),
	}

	forecasts <- currentCondition
	hourly := gjson.Get(jsonStr, "hourly")
	hourly.ForEach(func(key gjson.Result, value gjson.Result) bool {
		dtInt := value.Get("dt").Int()
		date := time.Unix(dtInt, 0)
		// fmt.Printf("Date: %v\n", dateStr)

		temp := value.Get("temp").Float()
		weatherId := gjson.Get(jsonStr, "weather.id").Int()

		// date, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", dateStr), location)

		hourlyForecast := forecast.HourlyForecast{
			TargetAt:    date,
			Temperature: float32(temp),
			Condition:   f.parseWeatherId(weatherId),
			Location:    city,
			Source:      f.source(),
		}

		forecasts <- hourlyForecast
		return true
	})

	daily := gjson.Get(jsonStr, "daily")
	daily.ForEach(func(key gjson.Result, value gjson.Result) bool {
		dtInt := value.Get("dt").Int()
		d := time.Unix(dtInt, 0)
		date := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, location)
		// fmt.Printf("Date: %v\n", dateStr)

		minTemp := value.Get("temp.min").Float()
		maxTemp := value.Get("temp.max").Float()
		weatherId := gjson.Get(jsonStr, "weather.id").Int()

		// date, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", dateStr), location)

		dailyForecast := forecast.DailyForecast{
			TargetAt:       date,
			TemperatureMin: float32(minTemp),
			TemperatureMax: float32(maxTemp),
			Condition:      f.parseWeatherId(weatherId),
			Location:       city,
			Source:         f.source(),
		}

		forecasts <- dailyForecast
		return true
	})
}

// 0: Sun, 1: Partial Clouds, 2: Clouds, 3: Fog, 4: Rain, 5: Snow
func (f *OpenWeatherFetcher) parseWeatherId(ID int64) int32 {
	if ID < 600 {
		return 4
	} else if ID < 700 {
		return 5
	} else if ID < 800 {
		return 3
	} else if ID == 800 {
		return 0
	} else if ID < 804 {
		return 1
	} else {
		return 2
	}
}
