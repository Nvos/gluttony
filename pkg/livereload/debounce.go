package livereload

import (
	"sync"
	"time"
)

func NewDebounce(after time.Duration) func(f func()) {
	d := &debounce{after: after, timer: nil, mu: sync.Mutex{}}

	return func(f func()) {
		d.add(f)
	}
}

type debounce struct {
	mu    sync.Mutex
	after time.Duration
	timer *time.Timer
}

func (d *debounce) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, f)
}
