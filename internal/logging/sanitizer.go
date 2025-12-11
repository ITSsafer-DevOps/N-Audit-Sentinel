package logging

import "regexp"

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// StripANSI removes ANSI escape sequences from a string
func StripANSI(s string) string {
	return ansi.ReplaceAllString(s, "")
}
