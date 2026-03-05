package livelog

import (
	"bytes"
	"strings"
	"sync"
)

// displayWriter is an io.Writer that routes complete lines through Display.Log.
type displayWriter struct {
	d       *Display
	mu      sync.Mutex
	partial bytes.Buffer
}

// Writer returns an io.Writer that buffers partial lines and dispatches
// complete lines (terminated by \n) through Display.Log.
// Safe for concurrent use by multiple goroutines.
func (d *Display) Writer() *displayWriter {
	return &displayWriter{d: d}
}

// Write implements io.Writer. Bytes are buffered until a newline is seen,
// then each complete line is dispatched to Display.Log.
func (w *displayWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.partial.Write(p)

	for {
		line, err := w.partial.ReadString('\n')
		if err != nil {
			// No newline found; put the partial data back
			w.partial.Reset()
			w.partial.WriteString(line)
			break
		}
		// Trim trailing newline since Log adds its own
		w.d.Log(strings.TrimRight(line, "\n"))
	}

	return len(p), nil
}
