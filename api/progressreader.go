package api

import (
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// refreshRate dictates how often ProgressReader writes to the progress channel.
const refreshRate = time.Millisecond * 200

// ProgressReader wraps a reader and reports its progress.
// A ProgressReader must be created by NewProgressReader.
type ProgressReader struct {
	io.Reader

	count int64 // all interaction on count must be guarded by atomic

	progress chan<- int64

	closeOnce sync.Once
	closed    chan struct{}
}

// NewProgressReader returns a new ProgressReader that wraps r and reports to p.
//
// The returned ProgressReader will not block sending to p.
// It is required to eventually Close() ProgressReader, which will also close p.
func NewProgressReader(r io.Reader, p chan<- int64) *ProgressReader {
	pr := &ProgressReader{
		Reader:   r,
		progress: p,
		closed:   make(chan struct{}),
	}

	go pr.start()

	return pr
}

// start will occasionally report the bytes read to the progress channel.
// Automatically fired in a new goroutine by NewProgressReader.
func (r *ProgressReader) start() {
	work := true
	last := int64(-1)

	for work {
		select {
		case <-r.closed:
			work = false
		case <-time.After(refreshRate):
			compare := atomic.LoadInt64(&r.count)

			// Only update on change.
			// Prevents sending updates every update if no further activity.
			if compare > last {
				r.update(compare)
				last = compare
			}
		}
	}

	// A last (still non-blocking) update to report final result, then close
	r.Update()
	close(r.progress)
}

// update triggers a non-blocking progress update
func (r *ProgressReader) update(x int64) {
	if r.progress != nil {
		select {
		case r.progress <- x:
		default:
		}
	}
}

// Update reports the bytes read to the progress channel.
//
// By default, ProgressReader will call Update periodically when more bytes have been written.
// Update is also called after Close, which will report the final total.
//
// Will not block. Manual calls to this function after Close() will panic.
func (r *ProgressReader) Update() {
	r.update(atomic.LoadInt64(&r.count))
}

// Read implements io.Reader.
func (r *ProgressReader) SetReader(newReader io.Reader) {
	r.Reader = newReader
}

// Read implements io.Reader.
func (r *ProgressReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)

	// Track bytes read
	atomic.AddInt64(&r.count, int64(n))

	return n, err
}

// Close implements io.Closer.
// Subsequent calls are ignored.
func (r *ProgressReader) Close() error {
	var err error

	r.closeOnce.Do(func() {

		// Close the underlying reader, if relevant
		if x, ok := r.Reader.(io.ReadCloser); ok {
			err = x.Close()
		}

		// signal the update loop
		close(r.closed)
	})

	return err
}
