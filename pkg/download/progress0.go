// progress.go

package download

import (
	"io"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

// ProgressTracker implements getter.ProgressTracker interface
// and uses cheggaaa/pb library for progress visualization.
type ProgressTracker0 struct {
	bars map[string]*pb.ProgressBar
}

// NewProgressTracker creates a new progress tracker instance.
func NewProgressTracker0() *ProgressTracker0 {
	return &ProgressTracker0{
		bars: make(map[string]*pb.ProgressBar),
	}
}

// TrackProgress implements getter.ProgressTracker interface.
func (t *ProgressTracker) TrackProgress0(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	filename := filepath.Base(src)
	bar := pb.Full.Start64(totalSize)
	bar.Set(pb.Bytes, true)

	bar.SetTemplateString(`{{string . "prefix" | green}} {{counters . }} {{bar . }} {{percent . }} {{speed . }}`)
	bar.Set("prefix", filename)

	// Store the bar in our map
	t.bars[src] = bar

	// Set the initial progress if we're resuming a download
	if currentSize > 0 {
		bar.SetCurrent(currentSize)
	}

	// Wrap the stream with our progress monitor
	reader := bar.NewProxyReader(stream)

	// Return a ReadCloser that will finish the progress bar on close
	return &progressReadCloser{
		Reader: reader,
		closer: func() error {
			bar.Finish()
			delete(t.bars, src)
			return stream.Close()
		},
	}
}

func (p *progressReadCloser) Close0() error {
	return p.closer()
}
