package webrtc

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pion/webrtc/v4"
)

type PeerID string

type Peers map[PeerID]*webrtc.PeerConnection

func (w *WebRTCServer) NewPeer(peerId PeerID, offer *webrtc.SessionDescription, trackerHandler TrackHandler, iceCandidateHandler OnNewICECandidateCallback) (*webrtc.PeerConnection, *webrtc.SessionDescription, error) {
	w.log.Info("inicializando novo peer")
	var err error
	// create new peer
	peerConn, err := w.api.NewPeerConnection(*w.PeerConfig)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateNewPeerConnection, err)
	}

	w.Peers[peerId] = peerConn

	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	} else if _, err = peerConn.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
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
	// Send the current time via a DataChannel to the remote peer every 3 seconds
	peerConn.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnOpen(func() {
			for range time.Tick(time.Second * 3) {
				if err = d.SendText(time.Now().String()); err != nil {
					panic(err)
				}
			}
		})
	})

	// peerConn.OnICECandidate(func(ice *webrtc.ICECandidate) {
	// 	if ice == nil {
	// 		return
	// 	}
	// 	iceCandidateHandler(ice)
	// })

	w.log.Info(fmt.Sprintf("setting remote description for client %s", peerId))
	err = peerConn.SetRemoteDescription(*offer)
	if err != nil {
		w.log.Error(fmt.Errorf("falha ao definir sdp remoto", err).Error())
	}

	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateAnswer, err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConn)

	err = peerConn.SetLocalDescription(answer)
	if err != nil {
		return nil, nil, multierror.Append(FailedToSetLocalDescription, err)
	}
	<-gatherComplete
	return peerConn, peerConn.LocalDescription(), nil
}
