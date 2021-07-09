package forecast

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

var conditions = [...]string{
	"Clear",
	"Partial Clouds",
	"Clouds",
	"Fog",
	"Rain",
	"Snow",
}

type HourlyForecast struct {
	gorm.Model
	TargetAt    time.Time
	Temperature float32
	Condition   int32 // 0: Sun, 1: Partial Clouds, 2: Clouds, 3: Rain, 4: Snow
	Location    string
	Source      string
}

func (f *HourlyForecast) DescribeCondition() string {
	return conditions[f.Condition]
}

func (f *HourlyForecast) Describe() string {
	when := "will be"
	if f.TargetAt.Before(time.Now()) {
		when = "was"
	}
	return fmt.Sprintf("According to %v in %v on the %v there %v %v and %v˚", f.Source, f.Location, f.TargetAt.Format("2006-01-02 15:04 -0700"), when, f.DescribeCondition(), f.Temperature)
}

// type Conditioned interface {
// 	getCondition() int32
// }

type DailyForecast struct {
	gorm.Model
	TargetAt       time.Time
	TemperatureMin float32
	TemperatureMax float32
	Condition      int32 // 0: Sun, 1: Partial Clouds, 2: Clouds, 3: Rain, 4: Snow
	Location       string
	Source         string
}

func (f *DailyForecast) DescribeCondition() string {
	return conditions[f.Condition]
}

func (f *DailyForecast) Describe() string {
	when := "will be"
	if f.TargetAt.Before(time.Now()) {
		when = "supposed to be"
	}
	return fmt.Sprintf("According to %v in %v on the %v there %v %v and %v˚-%v˚", f.Source, f.Location, f.TargetAt.Format("2006-01-02"), when, f.DescribeCondition(), f.TemperatureMin, f.TemperatureMax)
}
