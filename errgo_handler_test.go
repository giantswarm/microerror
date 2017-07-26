package microerror

import "testing"

func TestErrgoHandlerInterface(t *testing.T) {
	// This will not complie if ErrgoHandler does not fulfill Handler
	// interface.
	var _ Handler = NewErrgoHandler()
}
