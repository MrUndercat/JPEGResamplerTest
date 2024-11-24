package controller

import (
	myerrors "JPEGResempler/api/errors"
	"JPEGResempler/api/request"
	"JPEGResempler/api/response"
	"JPEGResempler/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	port          int
	JPEGResampler *service.Service
}

func NewHTTPHandler(port int, JPEGResampler *service.Service) *HTTPHandler {
	return &HTTPHandler{
		port:          port,
		JPEGResampler: JPEGResampler,
	}
}

func (h *HTTPHandler) ListenAndServe() {
	r := gin.Default()
	r.POST("/images", func(c *gin.Context) {
		var req request.Request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, response.ResponseError{
				Code: 400,
				Body: struct {
					Error string `json:"error"`
				}{
					Error: myerrors.ErrInvalidImageFormat.Error(),
				},
			})
			return
		}

		// Проверяем ошибку перед вызовом Error() и разделяем логику
		err := h.JPEGResampler.ResampleAndSaveImage(req.Image)
		if err != nil {
			if err.Error() == "file already exists" {
				// Обрабатываем случай, когда файл уже существует
				c.JSON(200, response.ResponseOK{
					Code: 200,
					Body: struct {
						Time   int  `json:"time"`
						Cached bool `json:"cached"`
					}{
						Time:   0,
						Cached: true,
					},
				})
			} else {
				// Обрабатываем другие ошибки
				c.JSON(400, response.ResponseError{
					Code: 400,
					Body: struct {
						Error string `json:"error"`
					}{
						Error: err.Error(),
					},
				})
			}
		} else {
			// Успешная обработка
			c.JSON(200, response.ResponseOK{
				Code: 200,
				Body: struct {
					Time   int  `json:"time"`
					Cached bool `json:"cached"`
				}{
					Time:   0,
					Cached: false,
				},
			})
		}
	})
	r.Run(":" + fmt.Sprintf("%d", h.port))
}
