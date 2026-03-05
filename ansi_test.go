package livelog

import (
	"strings"
	"testing"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no escapes", "hello world", "hello world"},
		{"red text", "\033[31mred\033[0m", "red"},
		{"bold green", "\033[1;32mbold green\033[0m", "bold green"},
		{"cursor movement", "\033[2K\033[1Acursor", "cursor"},
		{"empty", "", ""},
		{"multiple sequences", "\033[1m\033[31mhello\033[0m \033[34mworld\033[0m", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripANSI(tt.in)
			if got != tt.want {
				t.Errorf("StripANSI(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestVisibleWidth(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"plain", "hello", 5},
		{"with color", "\033[31mhello\033[0m", 5},
		{"empty", "", 0},
		{"only ansi", "\033[31m\033[0m", 0},
		{"mixed", "ab\033[31mcd\033[0mef", 6},
		{"40 chars with color", "\033[31m" + strings.Repeat("X", 40) + "\033[0m", 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VisibleWidth(tt.in)
			if got != tt.want {
				t.Errorf("VisibleWidth(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		maxWidth int
		wantMax  int // visible width should be <= this
	}{
		{"no truncation needed", "hello", 10, 5},
		{"exact fit", "hello", 5, 5},
		{"truncate plain", "hello world", 5, 5},
		{"truncate with ansi", "\033[31m" + strings.Repeat("X", 40) + "\033[0m", 20, 20},
		{"zero width", "hello", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Truncate(tt.in, tt.maxWidth)
			gotWidth := VisibleWidth(got)
			if gotWidth > tt.wantMax {
				t.Errorf("Truncate(%q, %d): visible width = %d, want <= %d",
					tt.in, tt.maxWidth, gotWidth, tt.wantMax)
			}
		})
	}
}

func TestTruncatePreservesAnsiReset(t *testing.T) {
	in := "\033[31m" + strings.Repeat("X", 40) + "\033[0m"
	got := Truncate(in, 10)
	if !strings.Contains(got, "\033[0m") {
		t.Errorf("Truncate should append ANSI reset, got %q", got)
	}
}

func TestTruncateNoTruncation(t *testing.T) {
	in := "short"
	got := Truncate(in, 20)
	if got != in {
		t.Errorf("Truncate(%q, 20) = %q, want %q (unchanged)", in, got, in)
	}
}
