package main

import (
	"doctor_recorder/internal/infrastructure"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	config, err := infrastructure.NewAppConfig()
	a, err := infrastructure.NewApp(config)
	if err != nil {
		panic(err)
	}
	for _, r := range a.Routes() {
		log.Printf("%s %s %s", r.Method, r.Path, r.Name)

	}
	a.Logger.Fatal(a.Start(":3000"))
}

func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "pages/home.tmpl", "World")
}
