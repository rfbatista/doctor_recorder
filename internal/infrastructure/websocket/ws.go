package websocket

import (
	"doctor_recorder/internal/infrastructure/logger"
	"doctor_recorder/internal/infrastructure/webrtc"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	pion "github.com/pion/webrtc/v3"
)

const (
	// time to read the next client's pong message
	pongWait = 60 * time.Second
	// time period to send pings to client
	pingPeriod = (pongWait * 9) / 10
	// time allowed to write a message to client
	writeWait = 10 * time.Second
	// max message size allowed
	maxMessageSize = 8192
	// I/O read buffer size
	readBufferSize = 4096
	// I/O write buffer size
	writeBufferSize = 4096
)

func NewWebsocket(log *logger.Logger, wrtc *webrtc.WebRTCServer) (*Websocket, error) {
	serv := make(ServerSubscriptions)
	client := make(ClientSubscriptions)
	serv[TopicWebRTC] = make(ServerTopicSubscribers)
	client[TopicWebRTC] = make(ClientTopicSubscribers)
	send := make(chan Message, 1)
	return &Websocket{log: log, webrtc: wrtc, SendChannel: send, ServerSubscriptions: serv, ClientSubscriptions: client}, nil
}

var upgrader = websocket.Upgrader{}

type Websocket struct {
	log                 *logger.Logger
	ws                  *websocket.Conn
	webrtc              *webrtc.WebRTCServer
	conn                *pion.PeerConnection
	SendChannel         chan Message
	ReceiveChannel      chan Message
	ServerSubscriptions ServerSubscriptions
	ClientSubscriptions ClientSubscriptions
}

func (w *Websocket) Handler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	//TODO: adicionar remocao do cliente de todos os topicos
	defer conn.Close()
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	clientId := SubscriberId(uuid.New().String())
	w.ClientSubscriptions[TopicId(TopicWebRTC)][clientId] = conn
	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			w.log.Error(fmt.Errorf("failed to read websocket message", err).Error())
			break
		}
		w.log.Info(fmt.Sprintf("reading message from client: %s", clientId))
		w.ProcessMessage(conn, clientId, msg)
	}
	return nil
}
