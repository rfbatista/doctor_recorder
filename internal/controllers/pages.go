package controllers

import (
	"doctor_recorder/internal/entities"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (c Controller) HomePage() echo.HandlerFunc {
	t := time.Now()
	version := entities.Version{Date: t.String()}
	return func(c echo.Context) error {
		return c.Render(http.StatusOK, "home.html", version)
	}
}
