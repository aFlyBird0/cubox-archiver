package util

// Trunc the string to N, use rune
func Trunc(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return string([]rune(s)[:length])
}
