package webrtc

import (
	"fmt"

	"github.com/pion/webrtc/v4"
)

func (w *WebRTCServer) AddIceCandidate(peerId PeerID, iceCandidate *webrtc.ICECandidateInit) {
	if _, exist := w.Peers[peerId]; exist {
		w.log.Info(fmt.Sprintf("adding new ice candidate to peer %s", peerId))
		conn := w.Peers[peerId]
		err := conn.AddICECandidate(*iceCandidate)
		if err != nil {
			w.log.Error(fmt.Errorf("failed to add ice candidate", err).Error())
		}
	}
}
