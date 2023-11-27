package transcriber

import (
	"doctor_recorder/internal/infrastructure/webrtc"
	"doctor_recorder/pkg/ffmpeg"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	pion "github.com/pion/webrtc/v4"
)

var (
	sampleRate      = 48000
	channels        = 2
	bufferSize      = 5 * 1024 * 1024
	bufferTreshhold = 4 * 1024 * 1024
)

func (t *Transcriber) HandleTracker() webrtc.TrackHandler {
	return func(track *pion.TrackRemote, receiver *pion.RTPReceiver) {
		t.log.Info("track received!!")
		codec := track.Codec()
		if strings.EqualFold(codec.MimeType, pion.MimeTypeOpus) {
			t.HandleTrackerStream(track)
		}
	}
}

func (t *Transcriber) HandleTrackerStream(track *pion.TrackRemote) {
  pr, pw := io.Pipe()
	ff, err := ffmpeg.NewFfmpeg(t.log)
	if err != nil {
		t.log.Error("failed to create ffmpeg")
		return
	}
	trStream, err := NewWhisperStream(pr)
	if err != nil {
		t.log.Error("failed to create whisper stream")
		return
	}
  err = ff.SetOutput(pw)
	if err != nil {
    t.log.Error(fmt.Errorf("failed to create ffmpeg stream : %s", err).Error())
		return
	}
	oggFile, err := oggwriter.NewWith(ff, 48000, 2)
	if err != nil {
		t.log.Error("failed to create ogg writer")
		return
	}
	// opusdecoder, err := opusdecod.NewDecoder()
	// if err != nil {
	// 	t.log.Error("failed to create new opus decord")
	// 	return
	// }
	defer func() {
		ff.Close()
	}()
	audioStream := make(chan []byte)
	errs := make(chan error, 2)
	response := make(chan bool)
	timer := time.NewTimer(10 * time.Second)
	go func() {
		for {
			packet, _, err := track.ReadRTP()
			timer.Reset(1 * time.Second)
			if err != nil {
				timer.Stop()
				if err == io.EOF {
					close(audioStream)
					return
				}
				errs <- err
				return
			}
			if packet == nil {
				t.log.Warning("invalid nil packet")
				return
			}
			if len(packet.Payload) == 0 {
				return
			}
			oggFile.WriteRTP(packet)
			// audioStream <- packet.Payload
			// <-response
		}
	}()
	err = nil
  go func ()  {
    b := make([]byte, 512 * 512)
    for {
      t.log.Info("waiting to receive data")
      size, err := pr.Read(b)
      if err != nil {
        pr.Close()
        return
      }
      if size == 0 {
        continue
      }
      trStream.Write(b[:size])
      b = make([]byte, 512 * 512)
    }
  }()
	for {
		select {
		case _ = <-audioStream:
			// t.log.Info(fmt.Sprintf("bytes: %v", audioChunk))
			// opusPacket := codecs.OpusPacket{}
			// if _, err := opusPacket.Unmarshal(audioChunk); err != nil {
			// 	// Only handle Opus packets
			// 	t.log.Error("failed to decode opus package")
			// }
			// t.log.Info(fmt.Sprintf("bytes %v", audioChunk))
			// payload, err := opusdecoder.Decode(audioChunk)
			// if err != nil {
			// 	t.log.Error("failed to decode audio")
			// 	return
			// }
			// transformedPkg, err := ff.Write(payload)
			// if err != nil {
			// 	t.log.Error(fmt.Errorf("failed to write to ffmpeg : %s", err).Error())
			// 	return
			// }
			// _, err = trStream.Write(transformedPkg.Bytes())
			// if err != nil {
			// 	t.log.Error("failed to write ffmpeg to stream")
			// 	return
			// }
			response <- true
			// if err != nil {
			// 	t.log.Error(fmt.Errorf("failed to write to ffmpeg", err).Error())
			// 	return
			// }
		case <-timer.C:
      t.log.Error("failed to decode audio : timeout")
			return
		case err = <-errs:
			t.log.Error(fmt.Errorf("received a error", err).Error())
			return
		}
	}
}
