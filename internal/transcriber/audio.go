package transcriber

import (
	"doctor_recorder/internal/infrastructure/webrtc"
	"doctor_recorder/pkg/ffmpeg"
	"fmt"
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
			oggFile, err := OggWriterNew("output.ogg", 48000, 2)
			if err != nil {
				t.log.Error("failed to create new ffmpeg")
				return
			}
			ff.Init()
			defer func() {
				ff.Close()
				oggFile.Close()
			}()
			t.log.Info("Got Opus track, saving to disk as output.opus (48 kHz, 2 channels)")
			// buffer := make([]byte, bufferSize)
			// var p bytes.Buffer
			rtpBuf := make([]byte, 1400)
			for {
				t.log.Info("reading packages")
				packageSize, _, err := track.Read(rtpBuf)
				if err != nil {
					t.log.Error("failed to read track package")
					return
				}

				t.log.Info("writing packages")
				wrote, err := ff.Write(rtpBuf[:packageSize])
				if err != nil {
					t.log.Error(fmt.Errorf("failed to write track package", err).Error())
				}
				t.log.Info(fmt.Sprintf("wrote %s", wrote))
				// rtpPacket, _, err := track.ReadRTP()
				// if err != nil {
				// 	t.log.Error("failed to read track package")
				// 	return
				// }
				// if err := oggFile.WriteRTP(rtpPacket); err != nil {
				// 	t.log.Error("failed to write ogg file")
				// 	return
				// }
				// if len(buffer) >= 10000000000 {
				// 	t.log.Info(fmt.Sprintf("Buffer is nearly full. Current size: %d", len(buffer)))
				// 	result, err := ff.Write(buffer)
				// 	if err != nil {
				// 		t.log.Error("failed to create audio")
				// 		return
				// 	}
				// 	outputFile, err := os.Create("out.mp3") // create new file
				// 	if err != nil {
				// 		t.log.Error("failed to create audio file")
				// 		return
				// 	}
				// 	defer outputFile.Close()
				// 	_, err = outputFile.Write(result.Bytes()) // write result buffer to file
				// 	if err != nil {
				// 		t.log.Error("failed to write audio file")
				// 		return
				// 	}
				// }
			}
		}
		return
	}
}
