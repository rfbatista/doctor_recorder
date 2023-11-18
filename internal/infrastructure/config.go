package infrastructure

import "doctor_recorder/internal/infrastructure/webrtc"

type AppConfig struct {
	WebRTCConfig webrtc.WebRTCConfig
}

func NewAppConfig() (AppConfig, error) {
	return AppConfig{WebRTCConfig: webrtc.WebRTCConfig{StunUrl: "stun:stun.l.google.com:19302"}}, nil
}
