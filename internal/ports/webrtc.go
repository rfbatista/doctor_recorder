package ports

import "github.com/labstack/echo/v4"

type WebRTC interface {
	Init(e *echo.Echo) error
}
