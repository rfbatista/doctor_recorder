package whisper

import (
	"bytes"
	"doctor_recorder/internal/entities"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewWhisperProvider() WhisperProvider {
	return WhisperProvider{}
}

type WhisperProvider struct {
}

func (w *WhisperProvider) GetAudioTranscription(audioPath string) (*entities.Transcription, error) {
	url := "http://localhost:4000/api/v1/transcribe"
	method := "POST"

	values := map[string]string{"audio_path": audioPath}
	payload, err := json.Marshal(values)
	// 	payload := strings.NewReader(`{
	//     "audio_path": audio_path
	// }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var transcription entities.Transcription
	err = json.NewDecoder(res.Body).Decode(&transcription)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &transcription, err
}
