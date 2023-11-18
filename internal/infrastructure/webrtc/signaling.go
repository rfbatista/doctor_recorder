package webrtc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pion/webrtc/v4"
)

func (w *WebRTCServer) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(webrtc.SessionDescription)
		if err := c.Bind(u); err != nil {
			return err
		}
		_, sd, err := w.NewPeer(u)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, sd)
	}
}
