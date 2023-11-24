package websocket

import "github.com/gorilla/websocket"

type SubscriptionCallback func(clientId SubscriberId, message Message) error

type SubscriberId string

type ServerTopicSubscribers map[SubscriberId]SubscriptionCallback

type ClientTopicSubscribers map[SubscriberId]*websocket.Conn

type ServerSubscriptions map[TopicId]ServerTopicSubscribers

type ClientSubscriptions map[TopicId]ClientTopicSubscribers

func (s *Websocket) SubscribeServerCallback(callback SubscriptionCallback, subscriberId string, topic string) {
	// if topic exist, check the client map
	if _, exist := s.ServerSubscriptions[TopicId(topic)]; exist {
		client := s.ServerSubscriptions[TopicId(topic)]

		// if client already subbed, stop the process
		if _, subbed := client[SubscriberId(subscriberId)]; subbed {
			return
		}

		// if not subbed, add to client map
		client[SubscriberId(subscriberId)] = callback
		return
	}

	// if topic does not exist, create a new topic
	s.ServerSubscriptions[TopicId(topic)] = make(ServerTopicSubscribers)

	// add the client to the topic
	s.ServerSubscriptions[TopicId(topic)][SubscriberId(subscriberId)] = callback
}

func (s *Websocket) Subscribe(conn *websocket.Conn, clientId SubscriberId, topic TopicId) {
	// if topic exist, check the client map
	if _, exist := s.ServerSubscriptions[TopicId(topic)]; exist {
		client := s.ClientSubscriptions[TopicId(topic)]

		// if client already subbed, stop the process
		if _, subbed := client[SubscriberId(clientId)]; subbed {
			return
		}

		// if not subbed, add to client map
		client[SubscriberId(clientId)] = conn
		return
	}

	// if topic does not exist, create a new topic
	s.ClientSubscriptions[TopicId(topic)] = make(ClientTopicSubscribers)

	// add the client to the topic
	s.ClientSubscriptions[TopicId(topic)][SubscriberId(clientId)] = conn
}
