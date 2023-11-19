package websocket

import (
	"doctor_recorder/internal/infrastructure/logger"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func NewWebsocket(log *logger.Logger) (*Websocket, error) {
	return &Websocket{log: log}, nil
}

var upgrader = websocket.Upgrader{}

type Websocket struct {
	log *logger.Logger
	ws  *websocket.Conn
}

func (w *Websocket) Handler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}
		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		w.log.Info(fmt.Sprintf("new websocket message \n %s", msg))
	}
}

func (w *Websocket) Handler(c echo.Context) error {
 return nil
}
