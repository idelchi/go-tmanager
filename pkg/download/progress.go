package download

import (
	"io"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

// ProgressTracker implements getter.ProgressTracker interface
// and uses cheggaaa/pb library for progress visualization
type ProgressTracker struct {
	pool  *pb.Pool
	bars  map[string]*pb.ProgressBar
	mutex sync.Mutex
}

// NewProgressTracker creates a new progress tracker instance
func NewProgressTracker() *ProgressTracker {
	pool, _ := pb.StartPool()
	return &ProgressTracker{
		pool: pool,
		bars: make(map[string]*pb.ProgressBar),
	}
}

// TrackProgress implements getter.ProgressTracker interface
func (t *ProgressTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	filename := filepath.Base(src)
	bar := pb.Full.Start64(totalSize)
	bar.Set(pb.Bytes, true)

	bar.SetTemplateString(`{{string . "prefix" | green}} {{counters . }} {{bar . }} {{percent . }} {{speed . }}`)
	bar.Set("prefix", filename)

	// Add bar to pool and map
	t.pool.Add(bar)
	t.bars[src] = bar

	if currentSize > 0 {
		bar.SetCurrent(currentSize)
	}

	reader := bar.NewProxyReader(stream)

	return &progressReadCloser{
		Reader: reader,
		closer: func() error {
			t.mutex.Lock()
			defer t.mutex.Unlock()

			bar.Finish()
			delete(t.bars, src)

			// If this was the last bar, stop the pool
			if len(t.bars) == 0 {
				t.pool.Stop()
			}

			return stream.Close()
		},
	}
}

// progressReadCloser wraps a Reader with a custom closer function
type progressReadCloser struct {
	io.Reader
	closer func() error
}

func (p *progressReadCloser) Close() error {
	return p.closer()
}
