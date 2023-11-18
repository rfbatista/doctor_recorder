package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c Controller) HomePage() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", nil)
	}
}
