package customerrors

// DuplicateValueError is error, returns when value exists
type DuplicateValueError struct {
	Err error
}

// Unwrap method of Error interface
func (err *DuplicateValueError) Unwrap() error {
	return err.Err
}

// Error method of Error interface
func (err *DuplicateValueError) Error() string {
	return err.Err.Error()
}
