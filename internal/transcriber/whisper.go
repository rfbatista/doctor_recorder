package transcriber

import (
	"doctor_recorder/internal/infrastructure/websocket"
	"doctor_recorder/pkg/whisper"
	"fmt"
	"io"
	"os"
)

func NewWhisperStream(r io.Reader, filename string) (*WhisperStream, error) {
	f, err := os.Create("./audios/" + filename + ".wav")
	if err != nil {
		return &WhisperStream{}, err
	}
	pr := whisper.NewWhisperProvider()
	return &WhisperStream{f: f, r: r, pr: &pr}, nil
}

type WhisperStream struct {
	io.Writer
	r       io.Reader
	f       *os.File
	results chan Result
	pr      *whisper.WhisperProvider
}

func (w *WhisperStream) Write(buffer []byte) (int, error) {
	// log.Printf("whisper stream received audio %v", buffer)
	// log.Printf("whisper strem received audio")
	_, err := w.f.Write(buffer)
	if err != nil {
		fmt.Println("falha ao escrever no buffer")
	}
	return 0, nil
}

func (w *WhisperStream) GetTranscription(clientId websocket.SubscriberId, web *websocket.Websocket) {
	info, _ := w.f.Stat()
	w.f.Close()
	go func(web *websocket.Websocket, path string) {
		trans, _ := w.pr.GetAudioTranscription(path)
		message := websocket.Message{
			Type:          websocket.MessageTypeTranscription,
			Action:        websocket.Publish,
			Topic:         websocket.TopicWebRTC,
			Transcription: trans,
		}
		web.Publish(websocket.TopicWebRTC, clientId, message)
	}(web, "/home/renan/projetos/doctor_recorder/audios/"+info.Name())

}

func (w *WhisperStream) Result() <-chan Result {
	return w.results
}

func (w *WhisperStream) Close() error {
	w.f.Close()
	return nil
}
