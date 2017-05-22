package btcticker

import "time"

type fetchResult struct {
	name string
	key  string
	data *FetchData
	err  error
}

type subscribe struct {
	fetcher Fetcher
	fail    bool
	from    string
	to      string
	key     string
	updates chan fetchResult
	closing chan chan error
}

// close сloses the subscription, the s.updates channel, and returns the last error.
func (s *subscribe) close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

// loop periodically processes the elements, sends them to s.updates, and exits.
// Fetch asynchronously.
func (s *subscribe) loop() {
	var done chan fetchResult
	var err error

	for {
		var startFetch <-chan time.Time
		if done == nil {
			startFetch = time.After(time.Second)
		}
		select {
		case <-startFetch:
			done = make(chan fetchResult, 1)
			go func() {
				data, err := s.fetcher.Fetch(s.from, s.to)
				done <- fetchResult{name: s.fetcher.Name(), key: s.key, data: data, err: err}
			}()
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		case result := <-done:
			done = nil
			s.updates <- result
		}
	}
}
