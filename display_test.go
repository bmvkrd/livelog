package livelog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

// helper to capture Display output via a pipe
func capturedDisplay(isTTY bool) (*Display, *os.File, *os.File) {
	r, w, _ := os.Pipe()
	d := New(w, WithForceTTY(isTTY))
	return d, r, w
}

func readAll(r *os.File, w *os.File) string {
	w.Close()
	data, _ := io.ReadAll(r)
	return string(data)
}

func TestLog_NonTTY(t *testing.T) {
	d, r, w := capturedDisplay(false)
	d.Log("hello")
	d.Log("world")
	d.Flush()
	output := readAll(r, w)

	if !strings.Contains(output, "hello\n") {
		t.Errorf("expected 'hello\\n' in output, got %q", output)
	}
	if !strings.Contains(output, "world\n") {
		t.Errorf("expected 'world\\n' in output, got %q", output)
	}
	// No ANSI codes in non-TTY mode
	if strings.Contains(output, "\033[") {
		t.Errorf("non-TTY output should not contain ANSI codes, got %q", output)
	}
}

func TestLogf_NonTTY(t *testing.T) {
	d, r, w := capturedDisplay(false)
	d.Logf("count: %d", 42)
	output := readAll(r, w)

	if !strings.Contains(output, "count: 42\n") {
		t.Errorf("expected 'count: 42\\n', got %q", output)
	}
}

func TestSetLive_NonTTY(t *testing.T) {
	d, r, w := capturedDisplay(false)
	d.SetLive([]string{"progress: 50%"})
	d.Log("a log message")
	d.SetLive([]string{"progress: 75%"})
	d.Flush()
	output := readAll(r, w)

	if strings.Contains(output, "\033[") {
		t.Errorf("non-TTY output should not contain ANSI codes")
	}
	if !strings.Contains(output, "progress: 50%") {
		t.Errorf("expected 'progress: 50%%' in output")
	}
	if !strings.Contains(output, "a log message") {
		t.Errorf("expected 'a log message' in output")
	}
	if !strings.Contains(output, "progress: 75%") {
		t.Errorf("expected 'progress: 75%%' in output")
	}
}

func TestSetLive_TTY_ContainsANSI(t *testing.T) {
	d, r, w := capturedDisplay(true)
	d.SetLive([]string{"line1", "line2"})
	d.SetLive([]string{"line3", "line4"})
	d.Flush()
	output := readAll(r, w)

	// Second SetLive should move cursor up 2 lines
	if !strings.Contains(output, "\033[2A") {
		t.Errorf("TTY mode should contain cursor-up sequence, got %q", output)
	}
	// Should use clear-to-end-of-screen
	if !strings.Contains(output, "\033[J") {
		t.Errorf("TTY mode should contain clear-to-end sequence, got %q", output)
	}
}

func TestLog_WithLiveZone_TTY(t *testing.T) {
	d, r, w := capturedDisplay(true)
	d.SetLive([]string{"live1", "live2"})
	d.Log("interleaved log")
	d.Flush()
	output := readAll(r, w)

	// The log should erase live zone, write log, redraw live zone
	if !strings.Contains(output, "interleaved log") {
		t.Errorf("expected log message in output")
	}
	if !strings.Contains(output, "live1") {
		t.Errorf("expected live line to be redrawn after log")
	}
}

func TestClearLive_TTY(t *testing.T) {
	d, r, w := capturedDisplay(true)
	d.SetLive([]string{"live"})
	d.ClearLive()

	// After clear, logging should not produce cursor movement
	d.Log("after clear")
	output := readAll(r, w)

	// The output should contain "after clear" without preceding cursor-up
	// (since live zone was cleared before the log)
	if !strings.Contains(output, "after clear") {
		t.Errorf("expected 'after clear' in output")
	}
}

func TestSetLive_Empty_ClearsZone(t *testing.T) {
	d, r, w := capturedDisplay(true)
	d.SetLive([]string{"live"})
	d.SetLive(nil)
	d.Log("no live zone")
	d.Flush()
	output := readAll(r, w)

	if !strings.Contains(output, "no live zone") {
		t.Errorf("expected 'no live zone' in output")
	}
}

func TestConcurrent_LogAndSetLive(t *testing.T) {
	d, r, w := capturedDisplay(true)

	var wg sync.WaitGroup

	// 100 goroutines logging concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			d.Logf("log %d", n)
		}(i)
	}

	// Concurrent live updates
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			d.SetLive([]string{fmt.Sprintf("progress: %d%%", i*2)})
		}
	}()

	wg.Wait()
	d.Flush()
	output := readAll(r, w)

	// All 100 log lines should appear
	for i := 0; i < 100; i++ {
		expected := fmt.Sprintf("log %d", i)
		if !strings.Contains(output, expected) {
			t.Errorf("missing log line %q in output", expected)
		}
	}
}

func TestIsTTY(t *testing.T) {
	d, _, w := capturedDisplay(true)
	if !d.IsTTY() {
		t.Error("expected IsTTY=true with WithForceTTY(true)")
	}
	w.Close()

	d2, _, w2 := capturedDisplay(false)
	if d2.IsTTY() {
		t.Error("expected IsTTY=false with WithForceTTY(false)")
	}
	w2.Close()
}

func TestTerminalWidth_NonTTY(t *testing.T) {
	d, _, w := capturedDisplay(false)
	defer w.Close()
	if d.TerminalWidth() != 80 {
		t.Errorf("non-TTY should return 80, got %d", d.TerminalWidth())
	}
}
