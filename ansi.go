package livelog

import "regexp"

// ansiPattern matches ANSI escape sequences (CSI sequences).
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// StripANSI removes all ANSI escape sequences from s.
func StripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

// VisibleWidth returns the number of visible columns s occupies
// after stripping ANSI escape codes. Does not account for wide (CJK) characters.
func VisibleWidth(s string) int {
	return len([]rune(StripANSI(s)))
}

// Truncate shortens s to fit within maxWidth visible columns.
// ANSI escape sequences are preserved and don't count toward the width.
// If truncation occurs, an ellipsis is appended and ANSI is reset.
func Truncate(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}
	if VisibleWidth(s) <= maxWidth {
		return s
	}

	runes := []rune(s)
	var result []rune
	visible := 0

	i := 0
	for i < len(runes) {
		// Check for ANSI escape sequence
		if runes[i] == '\x1b' && i+1 < len(runes) && runes[i+1] == '[' {
			// Find end of ANSI sequence
			j := i + 2
			for j < len(runes) && !isAnsiTerminator(runes[j]) {
				j++
			}
			if j < len(runes) {
				j++ // include terminator
			}
			result = append(result, runes[i:j]...)
			i = j
			continue
		}

		if visible >= maxWidth-1 {
			result = append(result, '…')
			result = append(result, []rune("\033[0m")...)
			break
		}

		result = append(result, runes[i])
		visible++
		i++
	}

	return string(result)
}

func isAnsiTerminator(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}
