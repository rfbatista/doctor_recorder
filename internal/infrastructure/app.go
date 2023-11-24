package infrastructure

import (
	"doctor_recorder/internal/controllers"
	"doctor_recorder/internal/infrastructure/logger"
	"doctor_recorder/internal/infrastructure/webrtc"
	"doctor_recorder/internal/infrastructure/websocket"
	"doctor_recorder/internal/transcriber"
	"doctor_recorder/internal/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewApp(config AppConfig) (*echo.Echo, error) {
	log := logger.NewLogger("app", nil)
	wrtc, err := webrtc.NewWebRTCServer(config.WebRTCConfig, log)
	if err != nil {
		panic(err)
	}
	err = wrtc.Init()
	if err != nil {
		panic(err)
	}
	ws, err := websocket.NewWebsocket(log, &wrtc)
	if err != nil {
		panic(err)
	}
	ts, err := transcriber.NewTranscriber(log, ws, &wrtc)
	if err != nil {
		panic(err)
	}
	ts.Init()
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Static("/static", "assets")
	e.GET("/ws", ws.Handler)
	t, err := view.NewTemplateEngine()
	if err != nil {
		return nil, err
	}
	err = t.Load()
	if err != nil {
		return nil, err
	}
	e.Renderer = &t
	c, err := controllers.NewController()
	if err != nil {
		return nil, err
	}
	_, err = c.Load(e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
