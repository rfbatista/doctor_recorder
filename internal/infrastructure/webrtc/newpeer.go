package webrtc

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/pion/webrtc/v3"
)

func (w *WebRTCServer) NewPeer(offer *webrtc.SessionDescription) (*webrtc.PeerConnection, *webrtc.SessionDescription, error) {
	w.log.Info("inicializando novo peer")
	var err error
	// create new peer
  peerConn, err := w.api.NewPeerConnection(w.PeerConfig)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateNewPeerConnection, err)
	}
	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	} else if _, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}
	peerConn.OnTrack(w.OnTrack(peerConn))
	peerConn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		w.log.Info("ice connection state changed")
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Ctrl+C the remote client to stop the demo")
		} else if connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateClosed {
			fmt.Println("Done writing media files")
			// Gracefully shutdown the peer connection
			if closeErr := peerConn.Close(); closeErr != nil {
				w.log.Error(closeErr)
			}
		}
	})
	peerConn.OnICECandidate(func(ice *webrtc.ICECandidate) {
		w.log.Info("new ice candidate found")
	})
	peerConn.SetRemoteDescription(*offer)
	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateAnswer, err)
	}
  w.log.Info("starting to gather ice")
	gatherComplete := webrtc.GatheringCompletePromise(peerConn)
	err = peerConn.SetLocalDescription(answer)
	if err != nil {
		return nil, nil, multierror.Append(FailedToSetLocalDescription, err)
	}
	<-gatherComplete
  w.log.Info("ice gathering is completed")
	return peerConn, &answer, nil
}
