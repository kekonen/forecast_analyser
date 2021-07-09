package main

// package main

import (
	"fmt"
	"forecast_analyser/lib/forecast"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

const json = `{"name":{"first":"Janet","last":"Prichard"},"age":47}`

func main1() {

	city := "Berlin"
	Source := "wttr.io"
	jsonFile, err := os.Open("wttr.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	// var jsonStr string
	jsonStr, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
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

	// value := gjson.Get(json, "name.last")
	fmt.Printf("%v\n%v\n", currentCondition, currentCondition.DescribeCondition())

	weatherDays := gjson.Get(string(jsonStr), "weather")
	weatherDays.ForEach(func(key gjson.Result, value gjson.Result) bool {
		dateStr := value.Get("date")
		fmt.Printf("Date: %v\n", dateStr)

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

		fmt.Printf("daily : %v\n", dailyForecast)

		value.Get("hourly").ForEach(func(key gjson.Result, value gjson.Result) bool {
			timeStr := value.Get("time").String()
			if len(timeStr) < 2 {
				timeStr = "0000"
			} else if len(timeStr) < 4 {
				timeStr = fmt.Sprintf("0%s", timeStr)
			}

			tempC := value.Get("tempC").Float()
			weatherCode := value.Get("weatherCode").Int()

			fmt.Printf(" * %v: %v %v\n", timeStr, tempC, parseWeatherCode(weatherCode))

			fullDate, err := time.ParseInLocation("2006-01-02T1504", fmt.Sprintf("%sT%s", dateStr, timeStr), location)
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

			fmt.Printf("hourly : %v\n", hourlyForecast)

			return true
		})
		return true
	})

	// fmt.Printf("%v\n")

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
