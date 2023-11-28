package transcriber

import (
	"doctor_recorder/internal/infrastructure/webrtc"
	"doctor_recorder/internal/infrastructure/websocket"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/pion/rtp"
	pion "github.com/pion/webrtc/v4"
)

var (
	sampleRate      = 48000
	channels        = 2
	bufferSize      = 5 * 1024 * 1024
	bufferTreshhold = 4 * 1024 * 1024
)

func (t *Transcriber) HandleTracker(clientId websocket.SubscriberId) webrtc.TrackHandler {
	c := t.HandleTrackerStream(clientId)
	return func(track *pion.TrackRemote, receiver *pion.RTPReceiver) {
		t.log.Info("track received!!")
		codec := track.Codec()
		if strings.EqualFold(codec.MimeType, pion.MimeTypeOpus) {
			c(track)
		}
	}
}

func (t *Transcriber) HandleTrackerStream(clientId websocket.SubscriberId) func(track *pion.TrackRemote) {
	return func(track *pion.TrackRemote) {
		pr, pw := io.Pipe()
		audioStream := make(chan *rtp.Packet)
		rawBytes := make(chan []byte)
		errs := make(chan error, 2)
		// response := make(chan bool)
		timer := time.NewTimer(10 * time.Second)
		sampling := time.NewTimer(5 * time.Second)
		// oggFile, err := oggwriter.NewWith(ff, 48000, 2)
		encoder := NewEncoder(pr, pw, t.log)
		encoder.Reset(clientId, t.ws)
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
				encoder.Write(packet)
				// audioStream <- packet
				// <-response
			}
		}()
		go func() {
			b := make([]byte, 512*512)
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
				rawBytes <- b[:size]
				// trStream.Write(b[:size])
				b = make([]byte, 512*512)
			}
		}()
		var err error = nil
		for {
			select {
			// case packet := <-audioStream:
			// 	continue

			case data := <-rawBytes:
				encoder.WriteWhisper(data)
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
				continue
				// if err != nil {
				// 	t.log.Error(fmt.Errorf("failed to write to ffmpeg", err).Error())
				// 	return
				// }
			case <-sampling.C:
				sampling.Reset(5 * time.Second)
				t.log.Error("reseting encodeer")
				encoder.Reset(clientId, t.ws)
				// response <- true
				// err := encoder.Reset()
				// trStream.GetTranscription(clientId, t.ws)
				// oggFile, _ = oggwriter.NewWith(ff, 48000, 2)
				if err != nil {
					t.log.Error("failed to create ogg writer")
					return
				}
				continue
			case <-timer.C:
				t.log.Error("failed to decode audio : timeout")
				return
			case err = <-errs:
				t.log.Error(fmt.Errorf("received a error", err).Error())
				return
			}
		}
	}
}
