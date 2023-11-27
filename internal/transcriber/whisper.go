package transcriber

import (
	"io"
	"log"
	"os"
)

func NewWhisperStream(r io.Reader) (io.Writer, error) {
	f, err := os.Create("xxxxx.wav")
	if err != nil {
		return &WhisperStream{}, err
	}
	return &WhisperStream{f: f, r: r}, nil
}

type WhisperStream struct {
	r       io.Reader
	f       *os.File
	results chan Result
}

func (w *WhisperStream) Write(buffer []byte) (int, error) {
	log.Printf("whisper stream received audio %v", buffer)
	// log.Printf("whisper stream received audio")
	w.f.Write(buffer)
	return 0, nil
}

func (w *WhisperStream) Result() <-chan Result {
	return w.results
}

func (w *WhisperStream) Close() error {
	w.f.Close()
	return nil
}
