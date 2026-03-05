package livelog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestWriter_CompleteLine(t *testing.T) {
	d, r, w := capturedDisplay(false)
	writer := d.Writer()
	fmt.Fprintln(writer, "hello world")
	output := readAll(r, w)

	if !strings.Contains(output, "hello world\n") {
		t.Errorf("expected 'hello world\\n', got %q", output)
	}
}

func TestWriter_PartialLines(t *testing.T) {
	r, pipeW, _ := os.Pipe()
	d := New(pipeW, WithForceTTY(false))
	writer := d.Writer()

	// Write partial line, then complete it
	writer.Write([]byte("hel"))
	writer.Write([]byte("lo\n"))

	pipeW.Close()
	data, _ := io.ReadAll(r)
	output := string(data)

	if !strings.Contains(output, "hello\n") {
		t.Errorf("expected single 'hello\\n', got %q", output)
	}
	// Should not have "hel" as a separate line
	lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
	for _, line := range lines {
		if line == "hel" {
			t.Errorf("partial line 'hel' should not appear as separate output")
		}
	}
}

func TestWriter_MultipleLines(t *testing.T) {
	d, r, w := capturedDisplay(false)
	writer := d.Writer()
	fmt.Fprintf(writer, "line1\nline2\nline3\n")
	output := readAll(r, w)

	for _, expected := range []string{"line1", "line2", "line3"} {
		if !strings.Contains(output, expected) {
			t.Errorf("expected %q in output, got %q", expected, output)
		}
	}
}

func TestWriter_ConcurrentWrites(t *testing.T) {
	d, r, w := capturedDisplay(false)
	writer := d.Writer()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Fprintf(writer, "msg %d\n", n)
		}(i)
	}
	wg.Wait()
	output := readAll(r, w)

	for i := 0; i < 50; i++ {
		expected := fmt.Sprintf("msg %d", i)
		if !strings.Contains(output, expected) {
			t.Errorf("missing %q in output", expected)
		}
	}
}

func TestWriter_WithLiveZone(t *testing.T) {
	d, r, w := capturedDisplay(true)
	writer := d.Writer()

	d.SetLive([]string{"progress bar"})
	fmt.Fprintln(writer, "a warning message")
	d.Flush()
	output := readAll(r, w)

	if !strings.Contains(output, "a warning message") {
		t.Errorf("expected warning message in output")
	}
	if !strings.Contains(output, "progress bar") {
		t.Errorf("expected progress bar in output")
	}
}
