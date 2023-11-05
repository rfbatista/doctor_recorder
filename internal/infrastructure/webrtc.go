package infrastructure

import "github.com/pion/webrtc/v2"

type WebRTCServer struct {
	config AppConfig
	api    *webrtc.API
}

func (w *WebRTCServer) Init() {
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{w.config.WebRTCConfig.Url},
			},
		},
	}
	m := webrtc.MediaEngine{}
	m.RegisterCodec(webrtc.NewRTPOpusCodec(0, 0))
	w.api = webrtc.NewAPI(webrtc.WithMediaEngine(m))
}
