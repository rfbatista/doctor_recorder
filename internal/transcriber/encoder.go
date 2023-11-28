package transcriber

import (
	"doctor_recorder/internal/infrastructure/logger"
	"doctor_recorder/internal/infrastructure/websocket"
	"doctor_recorder/pkg/ffmpeg"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/pion/rtp"

	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
)

func NewEncoder(pr io.Reader, pw io.Writer, log *logger.Logger) Encoder {
	return Encoder{pr: pr, pw: pw, log: log}
}

type Encoder struct {
	w       *oggwriter.OggWriter
	log     *logger.Logger
	out     io.Writer
	pr      io.Reader
	pw      io.Writer
	whisper *WhisperStream
}

func (e *Encoder) Write(packet *rtp.Packet) error {
	e.w.WriteRTP(packet)
	return nil
}

func (e *Encoder) WriteWhisper(b []byte) error {
	e.whisper.Write(b)
	return nil
}

func (e *Encoder) Reset(clientId websocket.SubscriberId, web *websocket.Websocket) error {
	if e.whisper != nil {
		e.whisper.GetTranscription(clientId, web)
	}
	ff, err := ffmpeg.NewFfmpeg(e.log)
	if err != nil {
		e.log.Error("failed to create ffmpeg")
		return err
	}
	err = ff.SetOutput(e.pw)
	if err != nil {
		e.log.Error(fmt.Errorf("failed to create ffmpeg stream : %s", err).Error())
		return err
	}
	trStream, err := NewWhisperStream(e.pr, uuid.NewString())
	if err != nil {
		e.log.Error("failed to create whisper stream")
		return err
	}
	e.whisper = trStream
	w, _ := oggwriter.NewWith(ff, 48000, 2)
	e.w = w
	return nil
}
