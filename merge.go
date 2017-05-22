package btcticker

type merge struct {
	subs    []*subscribe
	updates chan fetchResult
	quit    chan struct{}
	errs    chan error
}

// close closes the merged subscription, the updates channel, and returns the last error.
func (m *merge) close() error {
	close(m.quit)
	for range m.subs {
		if err := <-m.errs; err != nil {
			return err
		}
	}
	close(m.updates)
	return nil
}
