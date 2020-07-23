package microerror

import (
	"errors"
	"fmt"
	"strings"
)

func createStackTrace(err *stackedError) []StackEntry {
	stack := []StackEntry{
		err.stackEntry,
	}

	underlying := err.underlying
	var sErr *stackedError
	for errors.As(underlying, &sErr) {
		stack = append([]StackEntry{sErr.stackEntry}, stack...)
		underlying = sErr.underlying
	}

	return stack
}

func formatStackEntry(entry StackEntry) string {
	return fmt.Sprintf("%s:%d", entry.File, entry.Line)
}

func formatStackTrace(trace []StackEntry) string {
	var builder strings.Builder
	builder.Grow(len(trace))

	for _, stack := range trace {
		builder.WriteString(formatStackEntry(stack))
	}

	return builder.String()
}
