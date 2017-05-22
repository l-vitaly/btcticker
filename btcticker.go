package btcticker

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	tm "github.com/buger/goterm"
)

type counterData struct {
	fails int
	count int
}

type view struct {
	name      string
	amount    float64
	timestamp int
	mux       sync.Mutex
}

type btcTicker struct {
	cfg        Config
	subs       []*subscribe
	subscribes *merge
	view       map[string]*view
	runMux     sync.Mutex
}

// NewBtcTicker
func NewBtcTicker(cfg Config) *btcTicker {
	return &btcTicker{cfg: cfg, view: map[string]*view{}}
}

// Start
func (bt *btcTicker) Start() {
	go func() {
		tm.Clear()
		for _, p := range bt.cfg.Fetchers {
			if len(p.Exchanges)%2 == 1 {
				panic("exchanges: odd argument count")
			}
			for i := 0; i < len(p.Exchanges); i += 2 {
				bt.subs = append(bt.subs, bt.subscribe(getFetcher(p.Name), p.Exchanges[i], p.Exchanges[i+1]))
			}
		}

		bt.runMux.Lock()
		bt.subscribes = bt.merge(bt.subs...)
		go func() {
			for fr := range bt.subscribes.updates {
				bt.setView(fr)
				bt.render()
			}
		}()
		bt.runMux.Unlock()
	}()
}

// Stop
func (bt *btcTicker) Stop() {
	bt.runMux.Lock()
	defer bt.runMux.Unlock()
	bt.subscribes.close()
}

// setView sets a new currency amount, choosing the best amount
func (bt *btcTicker) setView(fr fetchResult) {
	if _, ok := bt.view[fr.key]; !ok {
		bt.view[fr.key] = &view{}
	}

	v := bt.view[fr.key]
	v.mux.Lock()
	defer v.mux.Unlock()

	if fr.err != nil {
		if fr.name == v.name {
			v.amount = 0
			v.name = ""
			v.timestamp = 0
		}
		return
	}

	// if there are no values
	if v.amount == 0 && v.name == "" && v.timestamp == 0 {
		v.name = fr.name
		v.amount = fr.data.Amount
		v.timestamp = fr.data.Timestamp
		return
	}

	// if the amount is smaller and its update date is longer
	if fr.data.Amount < v.amount && fr.data.Timestamp >= v.timestamp {
		v.name = fr.name
		v.amount = fr.data.Amount
		v.timestamp = fr.data.Timestamp
		return
	}

	// if there is an update of the current currency
	if v.name == fr.name && fr.data.Timestamp >= v.timestamp {
		v.amount = fr.data.Amount
		v.timestamp = fr.data.Timestamp
		return
	}
}

// render display changes in the console
func (bt *btcTicker) render() {
	tm.MoveCursor(1, 1)

	var viewParts []string
	for key, vi := range bt.view {
		viewParts = append(viewParts, fmt.Sprintf("%s: %.2f", strings.ToUpper(key), vi.amount))
	}

	sort.Strings(viewParts)

	counter := bt.getCounter()

	var sourceParts []string
	for name, c := range counter {
		sourceParts = append(sourceParts, fmt.Sprintf("%s (%d of %d)", strings.ToUpper(name), c.count-c.fails, c.count))
	}

	sort.Strings(sourceParts)

	tm.Println(fmt.Sprintf("%s\n\nActive sources: %s", strings.Join(viewParts, "  "), strings.Join(sourceParts, "  ")))

	tm.Flush()
}

// getCounter calculate the number of active fetchers
func (bt *btcTicker) getCounter() map[string]*counterData {
	counter := map[string]*counterData{}
	for _, s := range bt.subs {
		if _, ok := counter[s.key]; !ok {
			counter[s.key] = &counterData{}
		}
		if s.fail {
			counter[s.key].fails++
		}
		counter[s.key].count++
	}
	return counter
}

// merge
func (bt *btcTicker) merge(subs ...*subscribe) *merge {
	m := &merge{
		subs:    subs,
		updates: make(chan fetchResult),
		quit:    make(chan struct{}),
		errs:    make(chan error),
	}
	for _, subVal := range subs {
		go func(s *subscribe) {
			for {
				var fr fetchResult
				select {
				case fr = <-s.updates:
				case <-m.quit:
					m.errs <- s.close()
					return
				}
				select {
				case m.updates <- fr:
				case <-m.quit:
					m.errs <- s.close()
					return
				}
			}
		}(subVal)
	}
	return m
}

// subscribe
func (bt *btcTicker) subscribe(fetcher Fetcher, from, to string) *subscribe {
	s := &subscribe{
		fetcher: fetcher,
		from:    from,
		to:      to,
		key:     from + "/" + to,
		updates: make(chan fetchResult),
		closing: make(chan chan error),
	}
	go s.loop()
	return s
}
