// Fetcher for https://btc-e.com
package coindesk

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/l-vitaly/btcticker"
	"github.com/l-vitaly/btcticker/httputil"
)

var (
	errParse            = errors.New("parse error")
	errCurrencyNotFound = errors.New("currency not found")
)

const fetcherName = "btce"
const url = "https://btc-e.com/api/3/ticker/%s_%s"

type fetcher struct {
}

// Name
func (p *fetcher) Name() string {
	return fetcherName
}

// Fetch
func (p *fetcher) Fetch(from, to string) (*btcticker.FetchData, error) {
	r := httputil.NewRequest(60*time.Second, 60*time.Second, 60*time.Second)
	status, data, err := r.SendJSON(fmt.Sprintf(url, from, to), "GET", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errParse
	}
	key := fmt.Sprintf("%s_%s", from, to)
	currencyData := map[string]interface{}{}
	if currencyVal, ok := data[key]; ok {
		if currencyVal, ok := currencyVal.(map[string]interface{}); ok {
			currencyData = currencyVal
		}
	}
	updated := 0
	if updatedVal, ok := currencyData["updated"].(float64); ok {
		updated = int(updatedVal)
	}

	if amount, ok := currencyData["last"].(float64); ok {
		return &btcticker.FetchData{Amount: amount, Timestamp: updated}, nil
	}
	return nil, errCurrencyNotFound
}

func init() {
	btcticker.RegisterFetcher(&fetcher{})
}
