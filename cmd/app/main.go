package main

import (
	mycontroller "JPEGResempler/internal/controller/http_handler"
	myService "JPEGResempler/internal/service"
	myMiddleware "JPEGResempler/internal/service/app/middleware"
	myStorage "JPEGResempler/internal/service/app/storage"

	"flag"
	"log/slog"
	"os"
)

func main() {
	srcDir := flag.String("path-orig", "/tmp/img_orig/", "Directory for original images")
	resDir := flag.String("path-res", "/tmp/img_res/", "Directory for resampled images")
	width := flag.Uint("width", 200, "Size of resampled images")
	height := flag.Uint("height", 200, "Size of resampled images")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	storage := myStorage.NewStorage(*srcDir, *resDir, *width, *height, logger)

	middleware := myMiddleware.NewMiddleware()

	service := myService.NewService(middleware, storage, logger)

	handler := mycontroller.NewHTTPHandler(8085, service)

	if err := handler.ListenAndServe(); err != nil {
		logger.Error(err.Error())
	}
}
