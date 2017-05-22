package btcticker

import "sync"

var (
	fetchersMu sync.RWMutex
	fetchers   = make(map[string]Fetcher)
)

// FetchData fetcher work result
type FetchData struct {
	Amount    float64
	Timestamp int
}

// Fetcher the interface
type Fetcher interface {
	Name() string
	Fetch(from, to string) (*FetchData, error)
}

// getFetcher
func getFetcher(name string) Fetcher {
	fetchersMu.RLock()
	defer fetchersMu.RUnlock()
	if _, ok := fetchers[name]; !ok {
		panic("parser: GetParser called parser not found " + name)
	}
	return fetchers[name]
}

// RegisterFetcher
func RegisterFetcher(fetcher Fetcher) {
	fetchersMu.Lock()
	defer fetchersMu.Unlock()
	if fetcher == nil {
		panic("fetcher: RegisterParser fetcher is nil")
	}
	if _, dup := fetchers[fetcher.Name()]; dup {
		panic("fetcher: RegisterParser called twice for fetcher " + fetcher.Name())
	}
	fetchers[fetcher.Name()] = fetcher
}
