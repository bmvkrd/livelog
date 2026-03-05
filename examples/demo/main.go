// Command demo shows livelog's Display and ProgressBar working together.
// Run with: go run ./examples/demo
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bmvkrd/livelog"
)

func main() {
	d := livelog.New(os.Stdout)
	defer d.Flush()

	steps := []string{
		"Connecting to server",
		"Fetching metadata",
		"Downloading assets",
		"Verifying checksums",
		"Installing packages",
	}

	pb := &livelog.ProgressBar{
		Total: float64(len(steps)),
		Width: 30,
	}

	for i, step := range steps {
		pb.Current = float64(i)
		d.SetLive([]string{
			fmt.Sprintf("  Step: %s", step),
			fmt.Sprintf("  %s", pb.String()),
		})
		time.Sleep(600 * time.Millisecond)
		d.Log(fmt.Sprintf("✓ %s", step))
	}

	pb.Current = pb.Total
	d.SetLive([]string{fmt.Sprintf("  %s", pb.String())})
	time.Sleep(300 * time.Millisecond)
	d.ClearLive()
	d.Log("Done.")
}
