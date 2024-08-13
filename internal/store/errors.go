package store

type DuplicateValueError struct {
	Err error
}

func (err *DuplicateValueError) Unwrap() error {
	return err.Err
}

func (err *DuplicateValueError) Error() string {
	return err.Err.Error()
}
