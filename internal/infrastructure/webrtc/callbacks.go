package webrtc

import "github.com/pion/webrtc/v3"

type OnNewICECandidateCallback func(sd *webrtc.SessionDescription) error
type OnSendOfferCallback func(sd *webrtc.SessionDescription) error

func (w *WebRTCServer) SetOnNewICECandidateCallback(callback OnNewICECandidateCallback) {
	w.onNewICECandidateCallback = callback
}

func (w *WebRTCServer) SetOnSendOfferCallback(callback OnSendOfferCallback) {
	w.onSendOfferCallback = callback
}
