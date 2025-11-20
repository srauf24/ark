package router

import (
	"ark/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/api/status", h.Health.CheckHealth)

	r.Static("/static", "static")

	r.GET("/docs", h.OpenAPI.ServeOpenAPIUI)
}
