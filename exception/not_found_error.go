package exception

type NotFoundError struct {
	Error string
}

// Buat objek NewNotFoundError (constructror)
func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}
