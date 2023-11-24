package webrtc

import (
	"doctor_recorder/internal/infrastructure/logger"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/samplebuilder"
)

var (
	audioBuilder, videoBuilder     *samplebuilder.SampleBuilder
	audioTimestamp, videoTimestamp time.Duration
	streamKey                      string
)

type WebRTCServer struct {
	config     WebRTCConfig
	api        *webrtc.API
	PeerConfig webrtc.Configuration
	log        *logger.Logger
	Peers      Peers
}

func NewWebRTCServer(c WebRTCConfig, log *logger.Logger) (WebRTCServer, error) {
	return WebRTCServer{config: c, log: log, Peers: make(Peers)}, nil
}

func (w *WebRTCServer) Init() error {
	m := &webrtc.MediaEngine{}
	err := setupCodecs(m)
	if err != nil {
		return multierror.Append(FailedToSetupCodecs, err)
	}

	// i, err := setupInterceptors(m)
	if err != nil {
		return multierror.Append(FailedToSetupInterceptors, err)
	}
	w.PeerConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{w.config.StunUrl},
			},
		},
	}
	w.api = webrtc.NewAPI(webrtc.WithMediaEngine(m))
	return nil
}

func (w *WebRTCServer) handleData(track *webrtc.TrackRemote) {
	for {
		_, _, err := track.ReadRTP()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func pushOpus(rtpPacket *rtp.Packet) {
	audioBuilder.Push(rtpPacket)

	for {
		sample := audioBuilder.Pop()
		if sample == nil {
			return
		}
		// if audioWriter != nil {
		// 	audioTimestamp += sample.Duration
		// 	if _, err := audioWriter.Write(true, int64(audioTimestamp/time.Millisecond), sample.Data); err != nil {
		// 		panic(err)
		// 	}
		// }
	}
}
