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

	// Remove the 'Error: ' prefix if it exists
	if strings.HasPrefix(strings.ToLower(message), strings.ToLower(prefix)) {
		message = message[len(prefix):]
	}
	// This suffix is usually present in microerrors
	// without annotations.
	message = strings.TrimSuffix(message, " error")

	{
		// Capitalize the first letter.
		tmpMessage := []rune(message)
		if len(tmpMessage) > 0 {
			tmpMessage[0] = unicode.ToUpper(tmpMessage[0])
		}
		message = string(tmpMessage)
	}

	message = prefix + message

	return message
}
