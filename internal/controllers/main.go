package controllers

import "github.com/labstack/echo/v4"

type Controller struct {
}

func NewController() (Controller, error) {
	return Controller{}, nil
}

func (c Controller) Load(e *echo.Echo) (*echo.Echo, error) {
	e.GET("/", c.HomePage())
	return e, nil
}
