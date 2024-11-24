package service

import (
	"log/slog"
)

type ServiceI interface {
	ResampleAndSaveImage(image string, resampleSize int) error
}

type MiddlewareI interface {
	ValidateImage(image string) error
}

type StorageI interface {
	ProcessAndSaveImage(image string) error
}

type Service struct {
	Middleware MiddlewareI
	Storage    StorageI
	logger     *slog.Logger
}

func NewService(middleware MiddlewareI, storage StorageI, logger *slog.Logger) *Service {
	return &Service{
		Middleware: middleware,
		Storage:    storage,
		logger:     logger,
	}
}

func (s *Service) ResampleAndSaveImage(image string) error {
	if err := s.Middleware.ValidateImage(image); err != nil {
		s.logger.Error(err.Error())
		return err
	}
	if err := s.Storage.ProcessAndSaveImage(image); err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}
