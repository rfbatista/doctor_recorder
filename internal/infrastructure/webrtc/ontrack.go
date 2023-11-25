package webrtc

import (
	"fmt"
	"strings"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
)

func (w *WebRTCServer) OnTrack(peerConnection *webrtc.PeerConnection) func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
	return func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		w.log.Info("received a new track")
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				errSend := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}})
				if errSend != nil {
					fmt.Println(errSend)
				}
			}
		}()

		// Read incoming RTCP packets
		// Before these packets are returned they are processed by interceptors. For things
		// like TWCC and RTCP Reports this needs to be called.
		go func() {
			rtcpBuf := make([]byte, 1500)
			for {
				if _, _, rtcpErr := receiver.Read(rtcpBuf); rtcpErr != nil {
					return
				}
			}
		}()

		codec := track.Codec()
		if strings.EqualFold(codec.MimeType, webrtc.MimeTypeOpus) {
			fmt.Println("nada ")
			// saveToDisk(oggFile, track)
		} else if strings.EqualFold(codec.MimeType, webrtc.MimeTypeVP8) {
			fmt.Println("Got VP8 track, saving to disk as output.ivf")
			// saveToDisk(ivfFile, track)
		}
	}
}
