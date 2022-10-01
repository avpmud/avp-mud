package input

import (
	"strings"
)

var replacer = strings.NewReplacer(
	"\r", "",
	"\n", "",
)

// Sanitize strips newline and carriage returns from input.
func Sanitize(input string) string {
	return replacer.Replace(input)
}
