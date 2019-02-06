package microerror

type MarshalableError struct {
	Message string `json:"message"`
}

func (m *MarshalableError) Error() string {
	return m.Message
}
