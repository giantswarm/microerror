package microerror

// DeprecatedSetStack overwrites the stack of the current error with the stack
// of the error given. In case the given src error does not carry any stack, the
// current stack is not overwritten. This functionality is considered to be
// temporary and should only be used where absolutely necessary. One example we
// care about right now is the redirection of service errors in endpoints. In
// the example scenario a masked error is received, and for legacy reasons the
// endpoint does not want to forward the received error type, but rather its
// own. The downside of this technique is that we would lose the stack carried
// with the received error. Using DeprecatedSetStack the desired error type can
// be used and filled with the stack transported by the received error. The
// example below also shows how to preserve the original error message in such
// legacy situations.
func DeprecatedSetStack(dst, src error) error {
	srcMasked, ok := src.(*maskedError)
	if !ok {
		return dst
	}

	dstMasked := newMaskedError(dst)
	dstMasked.Stack = srcMasked.Stack

	return dstMasked
}
