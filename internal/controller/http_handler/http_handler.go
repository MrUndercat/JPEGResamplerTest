package controller

import (
	myerrors "JPEGResempler/api/errors"
	"JPEGResempler/api/request"
	"JPEGResempler/api/response"
	"JPEGResempler/internal/service"
	errors2 "JPEGResempler/internal/service/errors"

	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	port          int
	JPEGResampler *service.Service
	server        *http.Server
}

func NewHTTPHandler(port int, JPEGResampler *service.Service) *HTTPHandler {
	r := gin.Default()
	handler := &HTTPHandler{
		port:          port,
		JPEGResampler: JPEGResampler,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: r,
		},
	}

	r.POST("/images", func(c *gin.Context) {
		var req request.Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.Error{
				Code: http.StatusBadRequest,
				Body: struct {
					Error string `json:"error"`
				}{
					Error: myerrors.ErrInvalidImageFormat.Error(),
				},
			})
			return
		}

		err, time := handler.JPEGResampler.ResampleAndSaveImage(req.Image)
		if err != nil {
			if errors.As(err, &errors2.FileAlreadyExistsErr{"file already exists"}) {
				c.JSON(http.StatusOK, response.OK{
					Code: http.StatusOK,
					Body: struct {
						Time   int64 `json:"time"`
						Cached bool  `json:"cached"`
					}{
						Time:   time,
						Cached: true,
					},
				})
			} else {
				c.JSON(http.StatusBadRequest, response.Error{
					Code: http.StatusBadRequest,
					Body: struct {
						Error string `json:"error"`
					}{
						Error: err.Error(),
					},
				})
			}
		} else {
			c.JSON(http.StatusOK, response.OK{
				Code: http.StatusOK,
				Body: struct {
					Time   int64 `json:"time"`
					Cached bool  `json:"cached"`
				}{
					Time:   time,
					Cached: false,
				},
			})
		}
	})
	return handler
}

func (h *HTTPHandler) ListenAndServe() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	<-sigChan
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		h.JPEGResampler.Logger.Error(err.Error())
		return errors.New("server shutdown error")
	}
	h.JPEGResampler.Logger.Info("server successfully shut down")
	return nil
}
