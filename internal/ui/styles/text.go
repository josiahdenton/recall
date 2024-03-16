package styles

import "strings"

// Summary will take n characters in a string and add ellipses at the end
func Summary(s string, n int) string {
	if n >= len(s) {
		return s
	}
	var b strings.Builder
	b.WriteString(s[:n])
	b.WriteString("...")
	return b.String()
}
