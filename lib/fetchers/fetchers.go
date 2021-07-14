package fetchers

type Fetcher interface {
	Fetch(cities []string, forecastChan chan interface{})
}
