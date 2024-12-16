package utils

import (
	"regexp"
)

// StripColor removes the color tags from a string
func StripColor(input string) string {
	// Regular expression to match [color] tags
	re := regexp.MustCompile(`\[[^\]]*\]`)
	return re.ReplaceAllString(input, "")
}
