package validation

import (
	"strings"
	"unicode"
)

func RemoveWhitespace(s *string) string {
	var sb strings.Builder

	for _, r := range *s {
		if !unicode.IsSpace(r) {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

func CheckLength(input *string, possibleLengths *[]int) bool {
	for _, length := range *possibleLengths {
		if len(*input) == length {
			return true
		}
	}
	return false
}
