package livelog

import (
	"strings"
	"testing"
)

func TestProgressBar_Zero(t *testing.T) {
	pb := &ProgressBar{Total: 100, Current: 0}
	s := pb.String()
	if !strings.Contains(s, "0%") {
		t.Errorf("expected 0%% in output, got %q", s)
	}
	if !strings.Contains(s, "0s / 100s") {
		t.Errorf("expected '0s / 100s', got %q", s)
	}
}

func TestProgressBar_Half(t *testing.T) {
	pb := &ProgressBar{Total: 60, Current: 30}
	s := pb.String()
	if !strings.Contains(s, "50%") {
		t.Errorf("expected 50%% in output, got %q", s)
	}
	if !strings.Contains(s, "30s / 60s") {
		t.Errorf("expected '30s / 60s', got %q", s)
	}
}

func TestProgressBar_Full(t *testing.T) {
	pb := &ProgressBar{Total: 40, Current: 40}
	s := pb.String()
	if !strings.Contains(s, "100%") {
		t.Errorf("expected 100%% in output, got %q", s)
	}
	// Should have no empty chars
	if strings.Contains(s, "░") {
		t.Errorf("full bar should have no empty chars, got %q", s)
	}
}

func TestProgressBar_Overflow(t *testing.T) {
	pb := &ProgressBar{Total: 40, Current: 50}
	s := pb.String()
	if !strings.Contains(s, "100%") {
		t.Errorf("overflow should clamp to 100%%, got %q", s)
	}
}

func TestProgressBar_CustomChars(t *testing.T) {
	pb := &ProgressBar{Total: 10, Current: 5, Width: 10, FilledChar: "#", EmptyChar: "-"}
	s := pb.String()
	if !strings.Contains(s, "#####-----") {
		t.Errorf("expected '#####-----', got %q", s)
	}
}

func TestProgressBar_DefaultWidth(t *testing.T) {
	pb := &ProgressBar{Total: 100, Current: 50}
	s := pb.String()
	// Default width is 40, so 20 filled + 20 empty
	filled := strings.Count(s, "█")
	empty := strings.Count(s, "░")
	if filled+empty != 40 {
		t.Errorf("default width should be 40, got %d filled + %d empty = %d", filled, empty, filled+empty)
	}
}

func TestProgressBar_SetRatio(t *testing.T) {
	pb := &ProgressBar{Total: 100, Width: 10}
	pb.SetRatio(0.75)
	if pb.Current != 75 {
		t.Errorf("SetRatio(0.75) with Total=100 should set Current=75, got %f", pb.Current)
	}
}

func TestProgressBar_SetRatio_Clamp(t *testing.T) {
	pb := &ProgressBar{Total: 100, Width: 10}

	pb.SetRatio(-0.5)
	if pb.Current != 0 {
		t.Errorf("SetRatio(-0.5) should clamp to 0, got %f", pb.Current)
	}

	pb.SetRatio(1.5)
	if pb.Current != 100 {
		t.Errorf("SetRatio(1.5) should clamp to 100, got %f", pb.Current)
	}
}

func TestProgressBar_ZeroTotal(t *testing.T) {
	pb := &ProgressBar{Total: 0, Current: 0}
	s := pb.String()
	if !strings.Contains(s, "0%") {
		t.Errorf("zero total should show 0%%, got %q", s)
	}
}
