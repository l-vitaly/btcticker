// Fetcher for https://spectrocoin.com/scapi
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

const fetcherName = "spectrocoin"
const url = "https://spectrocoin.com/scapi/ticker/%s/%s"

type fetcher struct {
}

// Name
func (p *fetcher) Name() string {
	return fetcherName
}

// Fetch
func (p *fetcher) Fetch(from, to string) (*btcticker.FetchData, error) {
	r := httputil.NewRequest(60*time.Second, 60*time.Second, 60*time.Second)
	status, data, err := r.SendJSON(fmt.Sprintf(url, to, from), "GET", nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errParse
	}
	to = strings.ToUpper(to)
	if amount, ok := data["last"].(float64); ok {
		if timestamp, ok := data["timestamp"].(float64); ok {
			return &btcticker.FetchData{Amount: amount, Timestamp: int(timestamp)}, nil
		}
	}
	return nil, errCurrencyNotFound
}

func init() {
	btcticker.RegisterFetcher(&fetcher{})
}
