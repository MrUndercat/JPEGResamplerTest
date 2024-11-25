package middleware

import (
	"encoding/base64"
	"net/http"

	"JPEGResempler/api/errors"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) ValidateImage(image string) error {
	imageData, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return errors.ErrInvalidImageFormat // Ошибка, если не удается декодировать строку
	}

	base64Len := len(imageData)
	if base64Len < 1024 || base64Len > 10*1024*1024 {
		return errors.ErrInvalidImageSize
	}

	imageType := http.DetectContentType(imageData)
	if imageType != "image/jpeg" {
		return errors.ErrInvalidImageFormat
	}

	return nil
}
