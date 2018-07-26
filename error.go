package microerror

var invalidConfigError = New("invalid config")

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return Cause(err) == invalidConfigError
}

var metadataNotFoundError = New("metadata not found")

// IsMetadataNotFound asserts metadataNotFoundError.
func IsMetadataNotFound(err error) bool {
	return Cause(err) == metadataNotFoundError
}

var wrongTypeError = New("wrong type")

// IsWrongTypeError asserts wrongTypeError.
func IsWrongTypeError(err error) bool {
	return Cause(err) == wrongTypeError
}
