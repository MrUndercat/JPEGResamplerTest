package service

import (
	"log/slog"
	"time"
)

type ServiceI interface {
	ResampleAndSaveImage(image string) (error, int64)
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
	Logger     *slog.Logger
}

func NewService(middleware MiddlewareI, storage StorageI, logger *slog.Logger) *Service {
	return &Service{
		Middleware: middleware,
		Storage:    storage,
		Logger:     logger,
	}
}

func (s *Service) ResampleAndSaveImage(image string) (error, int64) {
	now := time.Now()
	if err := s.Middleware.ValidateImage(image); err != nil {
		s.Logger.Error(err.Error())
		return err, time.Since(now).Milliseconds()
	}
	if err := s.Storage.ProcessAndSaveImage(image); err != nil {
		s.Logger.Error(err.Error())
		return err, time.Since(now).Milliseconds()
	}
	return nil, time.Since(now).Milliseconds()
}
