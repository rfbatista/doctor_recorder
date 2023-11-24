package websocket

import "fmt"

func (s *Websocket) Publish(topic TopicId, clientId SubscriberId, message Message) {
	// if topic does not exist, stop the process
	if _, exist := s.ClientSubscriptions[topic]; !exist {
		return
	}

	// if topic exist
	clients := s.ClientSubscriptions[topic]

	_, exists := clients[clientId]
	if exists {
		client := clients[clientId]
		s.log.Info(fmt.Sprintf("sending message for client: %s", clientId))
		s.Send(client, message)
	}
}

func (s *Websocket) PublishToServerSubscribers(topic TopicId, clientId SubscriberId, message Message) {
	// if topic does not exist, stop the process
	if _, exist := s.ServerSubscriptions[topic]; !exist {
		return
	}

	// if topic exist
	clients := s.ServerSubscriptions[topic]

	for id, serverClient := range clients {
		s.log.Info(fmt.Sprintf("sending message for server client: %s", id))
		serverClient(clientId, message)
	}
}
