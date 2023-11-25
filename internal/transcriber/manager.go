package transcriber

import (
	"doctor_recorder/internal/infrastructure/logger"
	"doctor_recorder/internal/infrastructure/webrtc"
	"doctor_recorder/internal/infrastructure/websocket"
	"encoding/base64"
	"encoding/json"
	"fmt"

	pion "github.com/pion/webrtc/v4"
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
	t.webrtc.AddIceCandidate(webrtc.PeerID(clientId), message.Ice)
}

func (t *Transcriber) SendIceCandidate(clientId websocket.SubscriberId, iceCandidate *pion.ICECandidate) {
	t.log.Info("sending new ice candidate")
	cand := iceCandidate.ToJSON()
	message := websocket.Message{
		Type:   websocket.MessageTypeIceCandidate,
		Action: websocket.Publish,
		Topic:  websocket.TopicWebRTC,
		Ice:    &cand,
	}
	t.ws.Publish(websocket.TopicWebRTC, websocket.SubscriberId(clientId), message)
}

func (t *Transcriber) ReceiveSdp(clientId websocket.SubscriberId, message websocket.Message) {
	var sdp *pion.SessionDescription
	t.Decode(message.Sdp, &sdp)
	_, sdp, err := t.webrtc.NewPeer(webrtc.PeerID(clientId), sdp, t.HandleTracker(), t.ServerNewIceCandidate(clientId))
	if err != nil {
		t.log.Error(fmt.Errorf("failed to add new peer in webrtc: %s", err).Error())
		return
	}
	// send sdp create for the new connection
	t.SendSdp(webrtc.PeerID(clientId), sdp)
}

func (t *Transcriber) SendSdp(clientId webrtc.PeerID, sdp *pion.SessionDescription) {
	t.log.Info("sending ice candidate")

	s, err := t.Encode(sdp)
	if err != nil {
		t.log.Error("falha ao codificar sdp")
		return
	}
	message := websocket.Message{
		Type:   websocket.MessageTypeSDP,
		Action: websocket.Publish,
		Topic:  websocket.TopicWebRTC,
		Sdp:    s,
	}
	t.ws.Publish(websocket.TopicWebRTC, websocket.SubscriberId(clientId), message)
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

func (t *Transcriber) Encode(obj interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (t *Transcriber) Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}
	t.log.Info(string(b))

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}
