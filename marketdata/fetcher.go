package marketdata

type Fetcher interface {
	Fetch() (MarketData, error)
}
