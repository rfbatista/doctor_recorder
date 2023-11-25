package websocket

import pion "github.com/pion/webrtc/v4"

type Action string

const (
	Publish     Action = "publish"
	Subscribe          = "subscribe"
	Unsubscribe        = "unsubscribe"
)

type MessageType string

const (
	MessageTypeSDP          MessageType = "sdp"
	MessageTypeIceCandidate             = "ice"
	MessageTypeError                    = "error"
)

type TopicId string

const (
	TopicWebRTC TopicId = "webrtc"
)

type Message struct {
	Type   MessageType            `json:"type"`
	Action Action                 `json:"action"`
	Topic  TopicId                `json:"topic"`
	Sdp    string                 `json:"sdp"`
	Ice    *pion.ICECandidateInit `json:"ice"`
	Error  string                 `json:"error"`
}
