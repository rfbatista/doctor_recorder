package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

func (w *Websocket) Send(conn *websocket.Conn, message Message) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		w.log.Info("failed to marshal websocket message")
		return
	}
	// send simple message
	conn.WriteMessage(websocket.TextMessage, []byte(jsonMessage))
}

// SendWithWait sends message to the websocket client using wait group, allowing usage with goroutines
func (w *Websocket) SendWithWait(conn *websocket.Conn, message Message, wg *sync.WaitGroup) {
	// send simple message
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		w.log.Info("failed to marshal websocket message")
		wg.Done()
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte(jsonMessage))

	// set the task as done
	wg.Done()
}
