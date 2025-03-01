package helpers

import (
	"regexp"
	"strings"
)

// ToPlural 转换为复数形式
func ToPlural(word string) string {
	// Check if the word is already plural
	if strings.HasSuffix(word, "s") {
		return word
	}

	// Create a regular expression to match common singular noun patterns
	re := regexp.MustCompile("(ch|sh|s|x|z)$")

	// If the word ends with a consonant followed by "y", replace "y" with "ies"
	if strings.HasSuffix(word, "y") && !re.MatchString(word) {
		return word[:len(word)-1] + "ies"
	}

	// Otherwise, simply add an "s" to the end of the word
	return word + "s"
}
