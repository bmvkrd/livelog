package livelog

import (
	"fmt"
	"strings"
)

const defaultBarWidth = 40

// ProgressBar renders a configurable text-based progress bar.
type ProgressBar struct {
	Total      float64 // Total value (e.g., duration in seconds)
	Current    float64 // Current value
	Width      int     // Bar width in characters (default: 40)
	FilledChar string  // Character for filled portion (default: "█")
	EmptyChar  string  // Character for empty portion (default: "░")
}

// SetRatio sets Current from a 0.0–1.0 ratio (clamped).
func (pb *ProgressBar) SetRatio(ratio float64) {
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	pb.Current = ratio * pb.Total
}

// String renders the progress bar.
// Output format: "████████░░░░░░░░ 45%  18s / 40s"
func (pb *ProgressBar) String() string {
	width := pb.Width
	if width <= 0 {
		width = defaultBarWidth
	}

	filled := pb.FilledChar
	if filled == "" {
		filled = "█"
	}
	empty := pb.EmptyChar
	if empty == "" {
		empty = "░"
	}

	var ratio float64
	if pb.Total > 0 {
		ratio = pb.Current / pb.Total
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}
	}

	filledCount := int(ratio * float64(width))
	if filledCount > width {
		filledCount = width
	}
	emptyCount := width - filledCount

	bar := strings.Repeat(filled, filledCount) + strings.Repeat(empty, emptyCount)

	currentSec := int(pb.Current)
	totalSec := int(pb.Total)

	return fmt.Sprintf("%s %3.0f%%  %ds / %ds", bar, ratio*100, currentSec, totalSec)
}
