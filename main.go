package main

import (
	"time"
)

type Forecast struct {
	CollectedAt time.Time
	TargetAt    time.Time
	Type        string // "now", "hourly", "daily"
	Temperature float32
	Condition   string
	Location    string
	Source      string
}

type Fetcher interface {
	fetch()
}

type WttrFetcher struct {
}

func main() {
	cities := []string{
		"Berlin",
	}

}
