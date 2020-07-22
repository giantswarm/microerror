package microerror

import (
	"errors"
	"strings"
	"unicode"
)

const (
	prefix = "Error: "
	suffix = " error"
)

func Pretty(err error) string {
	var message string

	// Check if it's an annotated error.
	var aErr *annotatedError
	isAnnotatedErr := errors.As(err, &aErr)
	if isAnnotatedErr {
		message = aErr.annotation
	} else {
		// This is either an unmasked microerror, or
		// a simple 'errors.New()' error.
		message = err.Error()
	}

	message = prettifyErrorMessage(message)

	return message
}

func prettifyErrorMessage(message string) string {
	if len(message) < 1 {
		return message
	}

	r := []rune(message)
	// Pre-assign with the "Error: " prefix.
	m := []rune(prefix)

	var (
		nestedPos  int
		nestedChar rune
	)

	for pos := 0; pos < len(r); pos++ {
		if len(r) > len(prefix) && pos < len(prefix) {
			// Peek the next chars to see if the message
			// has the 'error: ' prefix.
			for nestedPos, nestedChar = range strings.ToLower(prefix) {
				if unicode.ToLower(r[pos+nestedPos]) != nestedChar {
					break
				}

				// Ignore the 'error: ' prefix if it exists.
				if nestedPos == len(prefix)-1 {
					pos += len(prefix)
					// Capitalize the first letter of the message.
					m = append(m, unicode.ToUpper(r[pos]))
					pos++
				}
			}
		}

		if len(r)+pos > len(suffix) && pos == len(r)-len(suffix) {
			// Peek the next chars to see if the message
			// has the ' error' suffix.
			for nestedPos, nestedChar = range strings.ToLower(suffix) {
				if unicode.ToLower(r[pos+nestedPos]) != nestedChar {
					break
				}

				// Ignore the ' error' suffix if it exists.
				// This suffix is usually present in microerrors
				// without annotations.
				if nestedPos == len(suffix)-1 {
					pos += len(suffix)
				}
			}
		}

		// We're already done, let's exit.
		if pos == len(r) {
			break
		}

		char := r[pos]
		// Capitalize the first letter of the message.
		if pos == 0 {
			char = unicode.ToUpper(char)
		}
		m = append(m, char)
	}

	return string(m)
}
