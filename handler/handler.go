package handler

import (
	"net/http"

	"github.com/iw4p/url-shortener/internal/service"
	"github.com/labstack/echo"
)

type Handler struct {
	urlService *service.URLService
}

func NewHandler(urlService *service.URLService) *Handler {
	return &Handler{urlService: urlService}
}

type Response struct {
	OK      bool
	Message string
}

type URLRequest struct {
	Short    string `json:"short"`
	Original string `json:"original"`
}

type URLResponse struct {
	Short    string `json:"short"`
	Original string `json:"original"`
}

func (h *Handler) HealthCheck(c echo.Context) error {
	res := Response{
		OK:      true,
		Message: "heartbeat is ok",
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) ShortURL(c echo.Context) error {
	req := new(URLRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{OK: false, Message: "Invalid request"})
	}
	shortURL, err := h.urlService.GetShorten(c.Request().Context(), req.Original)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{OK: false, Message: "Failed to shorten URL"})
	}
	return c.JSON(http.StatusOK, shortURL)
}

func (h *Handler) Redirect(c echo.Context) error {
	short := c.Param("redirect")
	originalData, err := h.urlService.GetOriginal(c.Request().Context(), short)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{OK: false, Message: "Failed to retrieve original URL"})
	}

	return c.Redirect(http.StatusMovedPermanently, originalData.Original)
}

func (h *Handler) GetOriginalURL(c echo.Context) error {
	req := new(URLRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{OK: false, Message: "Invalid request"})
	}
	originalURL, err := h.urlService.GetOriginal(c.Request().Context(), req.Short)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{OK: false, Message: "Failed to retrieve original URL"})
	}

	res := &URLResponse{
		Original: originalURL.Original,
		Short:    req.Short,
	}
	return c.JSON(http.StatusOK, res)
}
