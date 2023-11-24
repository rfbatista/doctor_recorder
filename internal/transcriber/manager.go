package transcriber

import (
	"doctor_recorder/internal/infrastructure/logger"
	"doctor_recorder/internal/infrastructure/webrtc"
	"doctor_recorder/internal/infrastructure/websocket"
	"fmt"

	pion "github.com/pion/webrtc/v3"
)

func NewTranscriber(log *logger.Logger, ws *websocket.Websocket, webrtc *webrtc.WebRTCServer) (*Transcriber, error) {
	return &Transcriber{log: log, ws: ws, webrtc: webrtc}, nil
}

var (
	TranscriberID = "transcriber"
)

type Transcriber struct {
	log    *logger.Logger
	ws     *websocket.Websocket
	webrtc *webrtc.WebRTCServer
}

func (t *Transcriber) Init() {
	t.ws.SubscribeServerCallback(t.SubscriberCallback(), TranscriberID, string(websocket.TopicWebRTC))
}

func (t *Transcriber) SubscriberCallback() websocket.SubscriptionCallback {
	return func(clientId websocket.SubscriberId, message websocket.Message) error {
		switch message.Type {
		case websocket.MessageTypeSDP:
			t.ReceiveSdp(clientId, message)
		case websocket.MessageTypeIceCandidate:
			t.ReceiveIceCandidate(clientId, message)
		default:
			t.SendError(websocket.TopicWebRTC, clientId, websocket.ErrActionUnrecognizable)
		}
		return nil

	}
}

func (t *Transcriber) ReceiveIceCandidate(clientId websocket.SubscriberId, message websocket.Message) {
	t.webrtc.AddIceCandidate(webrtc.PeerID(clientId), &pion.ICECandidateInit{Candidate: message.Ice})
}

func (t *Transcriber) SendIceCandidate(clientId websocket.SubscriberId, iceCandidate *pion.ICECandidate) {
	t.log.Info("sending new ice candidate")
	cand := iceCandidate.ToJSON()
	message := websocket.Message{
		Type:   websocket.MessageTypeIceCandidate,
		Action: websocket.Publish,
		Topic:  websocket.TopicWebRTC,
		Ice:    cand.Candidate,
	}
	t.ws.Publish(websocket.TopicWebRTC, websocket.SubscriberId(clientId), message)
}

func (t *Transcriber) ReceiveSdp(clientId websocket.SubscriberId, message websocket.Message) {
	_, sdp, err := t.webrtc.NewPeer(webrtc.PeerID(clientId), message.Sdp, t.HandleTracker(), t.ServerNewIceCandidate(clientId))
	if err != nil {
		t.log.Error(fmt.Errorf("failed to add new peer in webrtc: %s", err).Error())
		return
	}
	// send sdp create for the new connection
	t.SendSdp(webrtc.PeerID(clientId), sdp)
}

func (t *Transcriber) SendSdp(clientId webrtc.PeerID, sdp *pion.SessionDescription) {
	t.log.Info("sending ice candidate")
	message := websocket.Message{
		Type:   websocket.MessageTypeSDP,
		Action: websocket.Publish,
		Topic:  websocket.TopicWebRTC,
		Sdp:    sdp,
	}
	t.ws.Publish(websocket.TopicWebRTC, websocket.SubscriberId(clientId), message)
}

func (t *Transcriber) HandleTracker() webrtc.TrackHandler {
	return func(track *pion.TrackRemote, receiver *pion.RTPReceiver) {
		t.log.Info("track received!!")
		return
	}
}

func (t *Transcriber) ServerNewIceCandidate(clientId websocket.SubscriberId) webrtc.OnNewICECandidateCallback {
	return func(ice *pion.ICECandidate) {
		if ice == nil {
			return
		}
		t.SendIceCandidate(clientId, ice)
		return
	}
}

func (t *Transcriber) SendError(topic websocket.TopicId, clientId websocket.SubscriberId, err string) {
	message := websocket.Message{
		Topic: websocket.TopicWebRTC,
		Type:  websocket.MessageTypeError,
		Error: err,
	}
	t.ws.Publish(websocket.TopicWebRTC, clientId, message)
}
