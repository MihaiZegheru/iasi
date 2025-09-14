package iasiutils

// Add any shared utility functions here, e.g. truncateString, error helpers, etc.

// truncateString returns the first n characters of s, appending ... if truncated
func TruncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
