package webrtc

import (
	"doctor_recorder/internal/infrastructure/logger"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pion/interceptor"
	"github.com/pion/interceptor/pkg/intervalpli"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"
)

var (
	audioBuilder, videoBuilder     *samplebuilder.SampleBuilder
	audioTimestamp, videoTimestamp time.Duration
	streamKey                      string
	peerConn                       *webrtc.PeerConnection
)

type WebRTCServer struct {
	config      WebRTCConfig
	api         *webrtc.API
	PeerConfig  webrtc.Configuration
	connections map[string]chan *webrtc.PeerConnection
	log         *logger.Logger
}

func NewWebRTCServer(c WebRTCConfig, log *logger.Logger) (WebRTCServer, error) {
	return WebRTCServer{config: c, log: log}, nil
}

func (w *WebRTCServer) Init() error {
	m := &webrtc.MediaEngine{}

	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus, ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        111,
	}, webrtc.RTPCodecTypeAudio); err != nil {
		return err
	}

	audioBuilder = samplebuilder.New(10, &codecs.OpusPacket{}, 48000)

	i := &interceptor.Registry{}
	intervalPliFactory, err := intervalpli.NewReceiverInterceptor()
	if err != nil {
		return err
	}
	if err = webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		return err
	}
	w.PeerConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{w.config.StunUrl},
			},
		},
	}
	i.Add(intervalPliFactory)
	w.api = webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))
	return nil
}

func (w *WebRTCServer) NewPeer(offer *webrtc.SessionDescription) (*webrtc.PeerConnection, *webrtc.SessionDescription, error) {
	w.log.Info("inicializando novo peer")
	var err error
	peerConn, err = w.api.NewPeerConnection(w.PeerConfig)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateNewPeerConnection, err)
	}
	peerConn.SetRemoteDescription(*offer)
	answer, err := peerConn.CreateAnswer(nil)
	if err != nil {
		return nil, nil, multierror.Append(FailedToCreateAnswer, err)
	}
	err = peerConn.SetLocalDescription(answer)
	if err != nil {
		return nil, nil, multierror.Append(FailedToSetLocalDescription, err)
	}
	peerConn.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		w.log.Info("recebendo track")
		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				rtcpSendErr := peerConn.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}})
				if rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()
		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), track.Codec().RTPCodecCapability.MimeType)
		for {
			// Read RTP packets being sent to Pion
			rtp, _, readErr := track.ReadRTP()
			if readErr != nil {
				if readErr == io.EOF {
					return
				}
				panic(readErr)
			}
			switch track.Kind() {
			case webrtc.RTPCodecTypeAudio:
				pushOpus(rtp)
			}
		}
	})
	peerConn.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		w.log.Info("ice connection state changed")
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Ctrl+C the remote client to stop the demo")
		} else if connectionState == webrtc.ICEConnectionStateFailed || connectionState == webrtc.ICEConnectionStateClosed {
			fmt.Println("Done writing media files")
			// Gracefully shutdown the peer connection
			if closeErr := peerConn.Close(); closeErr != nil {
				panic(closeErr)
			}
			os.Exit(0)
		}
	})
	gatherComplete := webrtc.GatheringCompletePromise(peerConn)
	<-gatherComplete
	return peerConn, &answer, nil
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
