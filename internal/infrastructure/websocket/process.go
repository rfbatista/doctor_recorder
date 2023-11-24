package websocket

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func NewErrorMessage(error string) Message {
	return Message{
		Type:  MessageTypeError,
		Error: error,
	}
}

func (w *Websocket) ProcessMessage(conn *websocket.Conn, clientID SubscriberId, msg []byte) error {
	// parse message
	m := Message{}
	if err := json.Unmarshal(msg, &m); err != nil {
		w.Send(conn, NewErrorMessage(ErrInvalidMessage))
	}

	// convert all action to lowercase and remove whitespace
	action := strings.TrimSpace(strings.ToLower(string(m.Action)))

	switch Action(action) {
	case Publish:
		w.log.Info(fmt.Sprintf("publishing message for client: %s", clientID))
		w.PublishToServerSubscribers(m.Topic, clientID, m)

	case Subscribe:
		w.Subscribe(conn, clientID, m.Topic)

	case Unsubscribe:
		w.Unsubscribe(clientID, m.Topic)

	default:
		w.Send(conn, NewErrorMessage(ErrActionUnrecognizable))
	}

	return nil
}
