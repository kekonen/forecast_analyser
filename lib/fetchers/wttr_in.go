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

var Source = "wttr.in"

type WttrFetcher struct {
}

func (f *WttrFetcher) Fetch(cities []string, dailyForecasts chan forecast.DailyForecast, hourlyForecasts chan forecast.HourlyForecast) {
	for _, city := range cities {
		if f.hasCity(&city) {
			f.fetchCity(city, dailyForecasts, hourlyForecasts)
		}
	}
}

func (f *WttrFetcher) hasCity(city *string) bool {
	return *city == "Berlin"
}

func (f *WttrFetcher) getCityUrl(city *string) string {
	return fmt.Sprintf("http://wttr.in/%s?format=j1", *city)
}

func (f *WttrFetcher) fetchCity(city string, dailyForecasts chan forecast.DailyForecast, hourlyForecasts chan forecast.HourlyForecast) {
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
		Source:      Source,
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
			Source:         Source,
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
				Source:      Source,
			}

			// fmt.Printf("%v\n", hourlyForecast)
			hourlyForecasts <- hourlyForecast

			return true
		})
		return true
	})
}

func parseWeatherCode(weatherCode int64) int32 {
	switch weatherCode {
	case 113:
		return 0
	case 116:
		return 1
	case 119:
		return 2
	case 122:
		return 2
	case 143:
		return 3
	case 176:
		return 4
	case 179:
		return 5
	case 182:
		return 4
	case 185:
		return 4
	case 200:
		return 4
	case 227:
		return 5
	case 230:
		return 4
	case 248:
		return 3
	case 260:
		return 3
	case 263:
		return 4
	case 266:
		return 4
	case 281:
		return 4
	case 284:
		return 4
	case 293:
		return 4
	case 296:
		return 4
	case 299:
		return 4
	case 302:
		return 4
	case 305:
		return 4
	case 308:
		return 4
	case 311:
		return 4
	case 314:
		return 4
	case 317:
		return 4
	case 320:
		return 4
	case 323:
		return 5
	case 326:
		return 5
	case 329:
		return 5
	case 332:
		return 5
	case 335:
		return 5
	case 338:
		return 5
	case 350:
		return 5
	case 353:
		return 4
	case 356:
		return 4
	case 359:
		return 4
	case 362:
		return 4
	case 365:
		return 4
	case 368:
		return 4
	case 371:
		return 4
	case 386:
		return 4
	case 389:
		return 4
	case 392:
		return 5
	case 395:
		return 5
	default:
		return 0
	}
}

// 113:                                    : Clear
// 113:                                    : Sunny
// 116:                                    : Partly cloudy
// 119:                                    : Cloudy
// 122:                                    : Overcast
// 143:                                    : Mist
// 176:                                    : Patchy rain possible
// 179:                                    : Patchy snow possible
// 182:                                    : Patchy sleet possible
// 185:                                    : Patchy freezing drizzle possible
// 200:                                    : Thundery outbreaks possible
// 227:                                    : Blowing snow
// 230:                                    : Blizzard
// 248:                                    : Fog
// 260:                                    : Freezing fog
// 263:                                    : Patchy light drizzle
// 266:                                    : Light drizzle
// 281:                                    : Freezing drizzle
// 284:                                    : Heavy freezing drizzle
// 293:                                    : Patchy light rain
// 296:                                    : Light rain
// 299:                                    : Moderate rain at times
// 302:                                    : Moderate rain
// 305:                                    : Heavy rain at times
// 308:                                    : Heavy rain
// 311:                                    : Light freezing rain
// 314:                                    : Moderate or heavy freezing rain
// 317:                                    : Light sleet
// 320:                                    : Moderate or heavy sleet
// 323:                                    : Patchy light snow
// 326:                                    : Light snow
// 329:                                    : Patchy moderate snow
// 332:                                    : Moderate snow
// 335:                                    : Patchy heavy snow
// 338:                                    : Heavy snow
// 350:                                    : Ice pellets
// 353:                                    : Light rain shower
// 356:                                    : Moderate or heavy rain shower
// 359:                                    : Torrential rain shower
// 362:                                    : Light sleet showers
// 365:                                    : Moderate or heavy sleet showers
// 368:                                    : Light snow showers
// 371:                                    : Moderate or heavy snow showers
// 386:                                    : Patchy light rain with thunder
// 389:                                    : Moderate or heavy rain with thunder
// 392:                                    : Patchy light snow with thunder
// 395:                                    : Moderate or heavy snow with thunder
