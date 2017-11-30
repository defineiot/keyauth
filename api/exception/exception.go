package exception

// APIException is openauth error
type APIException struct {
	msg  string
	Code int
}

func (e *APIException) Error() string {
	return e.msg
}

// NewAPIException is openauth api error
func NewAPIException(text string, code int) error {
	return &APIException{msg: text, Code: code}
}
