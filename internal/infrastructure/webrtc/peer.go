package webrtc

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pion/webrtc/v3"
)

type PeerID string

type Peers map[PeerID]*webrtc.PeerConnection

func (w *WebRTCServer) NewPeer(peerId PeerID, offer *webrtc.SessionDescription, trackerHandler TrackHandler, iceCandidateHandler OnNewICECandidateCallback) (*webrtc.PeerConnection, *webrtc.SessionDescription, error) {
	w.log.Info("inicializando novo peer")
	var err error
	// create new peer
	peerConn, err := w.api.NewPeerConnection(w.PeerConfig)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateNewPeerConnection, err)
	}
	w.Peers[peerId] = peerConn
	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		return nil, nil, err
	} else if _, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		return nil, nil, err
	}
	peerConn.OnTrack(trackerHandler)
	peerConn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		w.log.Info(fmt.Sprintf("connection with peer %s has changed: %s", peerId, connectionState.String()))
		if connectionState == webrtc.ICEConnectionStateConnected {
			w.log.Info(fmt.Sprintf("connection with peer %s has connected", peerId))
		} else if connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateClosed {
			w.log.Info(fmt.Sprintf("Done writing media files"))
			// Gracefully shutdown the peer connection
			if closeErr := peerConn.Close(); closeErr != nil {
				w.log.Error(closeErr.Error())
			}
		}
	})
	peerConn.OnICECandidate(func(ice *webrtc.ICECandidate) {
		iceCandidateHandler(ice)
	})
	peerConn.SetRemoteDescription(*offer)
	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateAnswer, err)
	}
	err = peerConn.SetLocalDescription(answer)
	if err != nil {
		return nil, nil, multierror.Append(FailedToSetLocalDescription, err)
	}
	return peerConn, &answer, nil
}
