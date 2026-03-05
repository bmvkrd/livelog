package livelog_test

import (
	"fmt"
	"os"

	"github.com/bmvkrd/livelog"
)

// ExampleDisplay_Log demonstrates writing plain log lines to a Display.
func ExampleDisplay_Log() {
	d := livelog.New(os.Stdout, livelog.WithForceTTY(false))
	d.Log("Starting task...")
	d.Logf("Processing item %d of %d", 1, 10)
	d.Flush()
	// Output:
	// Starting task...
	// Processing item 1 of 10
}

// ExampleDisplay_SetLive demonstrates the live zone interleaved with log output.
// In non-TTY mode every SetLive call simply prints the lines; in a real terminal
// the live zone is redrawn in place.
func ExampleDisplay_SetLive() {
	d := livelog.New(os.Stdout, livelog.WithForceTTY(false))
	d.SetLive([]string{"Progress: 50%"})
	d.Log("Completed step 1")
	d.SetLive([]string{"Progress: 100%"})
	d.Flush()
	// Output:
	// Progress: 50%
	// Completed step 1
	// Progress: 100%
}

// ExampleDisplay_Writer demonstrates routing an io.Writer through the Display log zone.
func ExampleDisplay_Writer() {
	d := livelog.New(os.Stdout, livelog.WithForceTTY(false))
	w := d.Writer()
	fmt.Fprintln(w, "line from io.Writer")
	// Output:
	// line from io.Writer
}

// ExampleProgressBar demonstrates rendering a text progress bar.
func ExampleProgressBar() {
	pb := &livelog.ProgressBar{
		Total:      60,
		Current:    30,
		Width:      10,
		FilledChar: "#",
		EmptyChar:  "-",
	}
	fmt.Println(pb.String())
	// Output:
	// #####-----  50%  30s / 60s
}

// ExampleStripANSI demonstrates removing ANSI escape codes from a string.
func ExampleStripANSI() {
	colored := "\033[31mred text\033[0m"
	fmt.Println(livelog.StripANSI(colored))
	// Output:
	// red text
}

// ExampleVisibleWidth demonstrates measuring the printable width of a string.
func ExampleVisibleWidth() {
	colored := "\033[1;32mbold\033[0m"
	fmt.Println(livelog.VisibleWidth(colored))
	// Output:
	// 4
}

// ExampleTruncate demonstrates truncating a string to a visible column limit.
// Truncate appends an ANSI reset after the ellipsis, so VisibleWidth reports
// only the printable characters.
func ExampleTruncate() {
	long := "Hello, World!"
	result := livelog.Truncate(long, 5)
	fmt.Println(livelog.VisibleWidth(result))
	// Output:
	// 5
}
