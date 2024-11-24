package errors

import "errors"

var (
	ErrInvalidImageFormat = errors.New("invalid format, should be jpg/jpeg")
	ErrInvalidImageSize   = errors.New("invalid size, should be between 1024B and 10MB")
	ErrorProcessingImage  = errors.New("error processing image")
	ErrInvalidRequestBody = errors.New("invalid request body")
)
