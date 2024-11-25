package errors

type FileAlreadyExistsErr struct {
	Message string
}

func (e FileAlreadyExistsErr) Error() string {
	return e.Message
}
