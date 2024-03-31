package main

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

func CheckLength(cardNumber *string, possibleLengths *[]int) bool {
	for _, length := range *possibleLengths {
		if len(*cardNumber) == length {
			return true
		}
	}
	return false
}
