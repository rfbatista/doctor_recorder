package entities

type Transcription struct {
	Result      []string `json:"result"`
	Language    string   `json:"language"`
	Probability float64  `json:"probability"`
}
