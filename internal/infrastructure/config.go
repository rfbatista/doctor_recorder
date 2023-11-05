package infrastructure

type AppConfig struct {
	WebRTCConfig WebRTCConfig
}

type WebRTCConfig struct {
	Url string
}
