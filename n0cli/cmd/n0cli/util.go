package main

import "strings"

// referenced hasMeta() in path/filepath
func hasMeta(s string) bool {
	magicChars := `*?[`
	return strings.ContainsAny(s, magicChars)
}
