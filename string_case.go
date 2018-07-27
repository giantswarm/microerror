package microerror

import (
	"strings"
	"unicode"
)

func toStringCase(input string) string {
	chunks := []string{}

	for i, s := range strings.Split(input, "") {
		r := []rune(s)
		if i != 0 && unicode.IsUpper(r[0]) {
			chunks = append(chunks, string(" "))
		}
		chunks = append(chunks, strings.ToLower(s))
	}

	return strings.Join(chunks, "")
}
