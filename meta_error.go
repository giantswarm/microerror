package microerror

import "github.com/juju/errgo"

type metaError struct {
	cause    error
	metadata map[string]string
}

// NewMetaError creates a new instance of metaError. Its purpose is to act as
// container for meta information associated to a specific error. Package
// internally metaError instances can be used independent from each other. The
// specific error type matching can be used as usual. Using the meta error looks
// as follows. In the beginning is a usual error defined on its own. This error
// is the root cause once emitted during runtime.
//
//     var notEnoughWorkersError = microerror.New("not enough workers")
//
// Emitting the error defined above looks as follows. In the described situation
// something happens to cause the error propagation. The defined error is given
// to NewMetaError in order to annotate it with the given context information.
// The meta error is masked as usual.
//
//     return microerror.Mask(microerror.NewMetaError(notEnoughWorkersError, map[string]string{
//         "available",   strconv.Itoa(info.Workers.CountPerCluster.Max),
//         "description", "The amount of requested guest cluster workers exceeds the available number of host cluster nodes.",
//         "kind",        "notEnoughWorkersError",
//         "ops-recipe",  "https://github.com/giantswarm/ops-recipes/blob/master/349-not-enough-workers.md",
//         "requested",   strconv.Itoa(len(request.Cluster.Workers)),
//     }))
//
// Down the stack contextual information can be retrieved from the masked error.
// In the described example the key kind results in the dispatched value
// notEnoughWorkersError.
//
//     val, err := microerror.FromMetaError(err, "kind")
//     if err != nil {
//         return microerror.Mask(err)
//     }
//
func NewMetaError(err error, meta map[string]string) error {
	if err == nil {
		return Maskf(invalidConfigError, "error must not be empty")
	}
	if err.Error() == "" {
		return Maskf(invalidConfigError, "%T.Error() must not be empty", err)
	}
	if meta == nil {
		return Maskf(invalidConfigError, "meta must not be empty")
	}

	m := &metaError{
		cause:    err,
		metadata: meta,
	}

	return m
}

func (m *metaError) Error() string {
	return m.cause.Error()
}

func FromMetaError(err error, key string) (string, error) {
	c := errgo.Cause(err)

	m, ok := c.(*metaError)
	if !ok {
		return "", Maskf(wrongTypeError, "expected '%T', got '%T'", &metaError{}, c)
	}

	v, ok := m.metadata[key]
	if !ok {
		return "", Maskf(metadataNotFoundError, "no value for key '%s'", key)
	}

	return v, nil
}

func metaErrorCause(err error) error {
	m, ok := err.(*metaError)
	if !ok {
		return err
	}

	return m.cause
}
