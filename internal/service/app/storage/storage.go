package storage

import (
	myErrors "JPEGResempler/api/errors"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type Storage struct {
	OriginalImageCacheDirectory  string
	ResampledImageCacheDirectory string
	Width                        uint
	Height                       uint
	Logger                       *slog.Logger
}

func NewStorage(originalImageCacheDirectory string, resampledImageCacheDirectory string, width uint, height uint, logger *slog.Logger) *Storage {
	return &Storage{
		OriginalImageCacheDirectory:  originalImageCacheDirectory,
		ResampledImageCacheDirectory: resampledImageCacheDirectory,
		Width:                        width,
		Height:                       height,
		Logger:                       logger,
	}
}

func (s *Storage) ProcessAndSaveImage(imageBase64 string) error {
	imageData, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		s.Logger.Error(err.Error())
		return myErrors.ErrInvalidImageFormat
	}
	originalHash := fmt.Sprintf("%x", sha256.Sum256([]byte(imageBase64)))
	resampleHash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s-%d-%d", imageBase64, s.Width, s.Height))))
	originalPath := filepath.Join(s.OriginalImageCacheDirectory, originalHash+".jpeg")
	resampledPath := filepath.Join(s.ResampledImageCacheDirectory, resampleHash+".jpeg")
	if _, err := os.Stat(originalPath); os.IsNotExist(err) {
		if err := SaveImage(originalPath, imageBase64); err != nil {
			s.Logger.Error(err.Error())
			return myErrors.ErrorProcessingImage
		}
	}
	if _, err := os.Stat(resampledPath); !os.IsNotExist(err) {
		return errors.New("file already exists")
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		s.Logger.Error(err.Error())
		return myErrors.ErrorProcessingImage
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		s.Logger.Error(err.Error())
		return myErrors.ErrorProcessingImage
	}

	m := resize.Resize(s.Width, s.Height, img, resize.Lanczos3)

	file, err := os.Create(resampledPath)
	if err != nil {
		s.Logger.Error(err.Error())
		return myErrors.ErrorProcessingImage
	}
	defer file.Close()
	err = jpeg.Encode(file, m, &jpeg.Options{
		Quality: 80,
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return myErrors.ErrorProcessingImage
	}

	return nil
}

func SaveImage(path string, image string) error {
	data, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %v", err)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil

}
