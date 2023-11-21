package websocket

import (
	"doctor_recorder/internal/infrastructure/logger"
	"doctor_recorder/internal/infrastructure/webrtc"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	pion "github.com/pion/webrtc/v3"
)

func NewWebsocket(log *logger.Logger, wrtc *webrtc.WebRTCServer) (*Websocket, error) {
	return &Websocket{log: log, webrtc: wrtc}, nil
}

var upgrader = websocket.Upgrader{}

type Message struct {
	Type string                   `json:"type"`
	Sdp  *pion.SessionDescription `json:"sdp"`
	Ice  *pion.ICECandidate       `json:"ice"`
}

type Websocket struct {
	log    *logger.Logger
	ws     *websocket.Conn
	webrtc *webrtc.WebRTCServer
	conn   *pion.PeerConnection
}

func (w *Websocket) Handler(c echo.Context) error {
	wsconn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	w.ws = wsconn
	if err != nil {
		return err
	}
	go func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			w.log.Info("waiting for new message")
			_, msg, err := ws.ReadMessage()
			if err != nil {
				c.Logger().Error(err)
			}
			if msg == nil {
				continue
			} else {
				var clientMessage Message
				err = json.Unmarshal(msg, &clientMessage)
				if err != nil {
					w.log.Error(err)
				} else {
					if clientMessage.Type == "sdp" {
						conn, sdp, _ := w.webrtc.NewPeer(clientMessage.Sdp)
						w.conn = conn
						sdpMessage := Message{
							Type: "sdp",
							Sdp:  sdp,
						}
						ws.WriteJSON(sdpMessage)
					}
					if clientMessage.Type == "ice" {
					}
				}
				w.log.Info(fmt.Sprintf("new websocket message %s", clientMessage.Type))
			}
		}
	}(wsconn)
	return nil
}
func (w *Websocket) Init() error {
	w.log.Info("defining callback")
	callback := func(ice *pion.ICECandidate) {
		iceMessage := Message{
			Type: "ice",
			Ice:  ice,
		}
		w.ws.WriteJSON(iceMessage)
	}
	w.webrtc.SetOnNewICECandidateCallback(callback)
	return nil
}
