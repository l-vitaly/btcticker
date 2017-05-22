// Fetcher for http://api.coindesk.com
package coindesk

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/l-vitaly/btcticker"
	"github.com/l-vitaly/btcticker/httputil"
)

var (
	errParse            = errors.New("parse error")
	errCurrencyNotFound = errors.New("currency not found")
)

const fetcherName = "coindesk"
const url = "http://api.coindesk.com/v1/bpi/currentprice/%s.json"

type fetcher struct {
}

// Name
func (p *fetcher) Name() string {
	return fetcherName
}

// Fetch
func (p *fetcher) Fetch(from, to string) (*btcticker.FetchData, error) {
	to = strings.ToUpper(to)
	r := httputil.NewRequest(60*time.Second, 60*time.Second, 60*time.Second)
	status, data, err := r.SendJSON(fmt.Sprintf(url, to), "GET", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errParse
	}
	currencyData := map[string]interface{}{}
	if bpi, ok := data["bpi"]; ok {
		if bpi, ok := bpi.(map[string]interface{}); ok {
			if currencyVal, ok := bpi[to].(map[string]interface{}); ok {
				currencyData = currencyVal
			}
		}
	}
	updated := time.Time{}
	if timeVal, ok := data["time"].(map[string]interface{}); ok {
		if updatedVal, ok := timeVal["updatedISO"].(string); ok {
			updated, err = time.Parse(time.RFC3339, updatedVal)
			if err != nil {
				return nil, err
			}
		}
	}
	if amount, ok := currencyData["rate_float"].(float64); ok {
		return &btcticker.FetchData{Amount: amount, Timestamp: updated.Second()}, nil
	}
	return nil, errCurrencyNotFound
}

func init() {
	btcticker.RegisterFetcher(&fetcher{})
}
