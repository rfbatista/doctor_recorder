package webrtc

import "github.com/pion/webrtc/v3"

type OnNewICECandidateCallback func(ice *webrtc.ICECandidate)
type OnSendOfferCallback func(sd *webrtc.SessionDescription)

func (w *WebRTCServer) SetOnNewICECandidateCallback(callback OnNewICECandidateCallback) {
	w.onNewICECandidateCallback = callback
}

func (w *WebRTCServer) SetOnSendOfferCallback(callback OnSendOfferCallback) {
	w.onSendOfferCallback = callback
}
