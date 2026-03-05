package livelog

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync"

	"golang.org/x/term"
)

// Display manages two terminal output zones: a scrolling log zone and a
// pinned live zone at the bottom. All methods are safe for concurrent use.
type Display struct {
	mu        sync.Mutex
	out       *os.File
	fd        int
	isTTY     bool
	forceTTY  *bool // non-nil means overridden
	separator bool  // blank line between log and live zones
	liveLines []string
	liveCount int // number of logical lines last rendered
	buf       bytes.Buffer
}

// Option configures a Display.
type Option func(*Display)

// WithSeparator controls whether a blank line is inserted between the
// scrolling log zone and the live zone. Defaults to false.
func WithSeparator(sep bool) Option {
	return func(d *Display) {
		d.separator = sep
	}
}

// WithForceTTY overrides automatic TTY detection.
func WithForceTTY(isTTY bool) Option {
	return func(d *Display) {
		d.forceTTY = &isTTY
		d.isTTY = isTTY
	}
}

// New creates a Display that writes to out.
// TTY detection is automatic; override with WithForceTTY for testing.
func New(out *os.File, opts ...Option) *Display {
	fd := int(out.Fd())
	d := &Display{
		out:   out,
		fd:    fd,
		isTTY: term.IsTerminal(fd),
		separator: true,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// IsTTY reports whether the Display is writing to a terminal.
func (d *Display) IsTTY() bool {
	return d.isTTY
}

// TerminalWidth returns the terminal width in columns, or 80 if unknown.
func (d *Display) TerminalWidth() int {
	if !d.isTTY {
		return 80
	}
	w, _, err := term.GetSize(d.fd)
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

// Log writes a line to the scrolling zone above the live region.
func (d *Display) Log(args ...any) {
	msg := fmt.Sprint(args...)
	d.logLine(msg)
}

// Logf writes a formatted line to the scrolling zone.
func (d *Display) Logf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	d.logLine(msg)
}

func (d *Display) logLine(msg string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isTTY {
		fmt.Fprintln(d.out, msg)
		return
	}

	d.buf.Reset()

	// Erase live zone
	if d.liveCount > 0 {
		d.buf.WriteString("\033[")
		d.buf.WriteString(strconv.Itoa(d.liveCount))
		d.buf.WriteString("A\r\033[J")
	}

	// Write log line
	d.buf.WriteString(msg)
	d.buf.WriteByte('\n')

	// Redraw live zone
	d.writeLiveLines()

	d.out.Write(d.buf.Bytes())
}

// SetLive replaces the live zone content. Each element becomes one line.
// Pass nil or empty slice to clear the live zone.
func (d *Display) SetLive(lines []string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isTTY {
		for _, line := range lines {
			fmt.Fprintln(d.out, line)
		}
		return
	}

	d.buf.Reset()

	// Erase previous live zone
	if d.liveCount > 0 {
		d.buf.WriteString("\033[")
		d.buf.WriteString(strconv.Itoa(d.liveCount))
		d.buf.WriteString("A\r\033[J")
	}

	// Update and render new live zone
	d.liveLines = make([]string, len(lines))
	copy(d.liveLines, lines)
	d.writeLiveLines()

	d.out.Write(d.buf.Bytes())
}

// ClearLive erases the live zone.
func (d *Display) ClearLive() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.isTTY || d.liveCount == 0 {
		d.liveLines = nil
		d.liveCount = 0
		return
	}

	d.buf.Reset()
	d.buf.WriteString("\033[")
	d.buf.WriteString(strconv.Itoa(d.liveCount))
	d.buf.WriteString("A\r\033[J")
	d.out.Write(d.buf.Bytes())

	d.liveLines = nil
	d.liveCount = 0
}

// Flush clears the live zone and should be called before program exit.
func (d *Display) Flush() {
	d.ClearLive()
}

// writeLiveLines writes current liveLines to buf and updates liveCount.
// Must be called with mu held.
func (d *Display) writeLiveLines() {
	tw := d.termWidthLocked()
	d.liveCount = len(d.liveLines)
	if d.separator && d.liveCount > 0 {
		d.buf.WriteByte('\n')
		d.liveCount++
	}
	for _, line := range d.liveLines {
		if tw > 0 && VisibleWidth(line) > tw-1 {
			line = Truncate(line, tw-1)
		}
		d.buf.WriteString(line)
		d.buf.WriteByte('\n')
	}
}

// termWidthLocked returns terminal width; must be called with mu held.
func (d *Display) termWidthLocked() int {
	if !d.isTTY {
		return 80
	}
	w, _, err := term.GetSize(d.fd)
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
