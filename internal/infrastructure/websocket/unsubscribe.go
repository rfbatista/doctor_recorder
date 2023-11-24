package websocket

// Unsubscribe removes a clients from a topic's client map
func (s *Websocket) Unsubscribe(clientID SubscriberId, topic TopicId) {
	// if topic exist, check the client map
	if _, exist := s.ClientSubscriptions[topic]; exist {
		client := s.ClientSubscriptions[topic]

		// remove the client from the topic's client map
		delete(client, clientID)
	}
}
