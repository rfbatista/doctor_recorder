package transcriber

import (
	"doctor_recorder/internal/infrastructure/webrtc"
	"doctor_recorder/pkg/ffmpeg"
	"fmt"
	"os"
	"strings"

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
			ff, err := ffmpeg.NewFfmpeg()
			if err != nil {
				t.log.Error("failed to create new ffmpeg")
				return
			}
      ff.Init()
			defer func() {
				ff.Close()
			}()
			t.log.Info("Got Opus track, saving to disk as output.opus (48 kHz, 2 channels)")
			buffer := make([]byte, bufferSize)
			// var p bytes.Buffer
			t.log.Info("reading packages")
			for {
				rtpPacket, _, err := track.ReadRTP()
				if err != nil {
					t.log.Error("failed to read track package")
					return
				}
				rtpPacket.MarshalTo(buffer)
				if len(buffer) >= bufferTreshhold {
					t.log.Info(fmt.Sprintf("Buffer is nearly full. Current size: %d", len(buffer)))
					result, err := ff.Read(buffer)
					if err != nil {
						t.log.Error("failed to create audio")
						return
					}
					outputFile, err := os.Create("out.mp3") // create new file
					if err != nil {
						t.log.Error("failed to create audio file")
						return
					}
					defer outputFile.Close()
					_, err = outputFile.Write(result.Bytes()) // write result buffer to file
					if err != nil {
						t.log.Error("failed to write audio file")
						return
					}
				}
			}
		}
		return
	}
}
