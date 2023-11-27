package ffmpeg

import (
	"bytes"
	"doctor_recorder/internal/infrastructure/logger"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func NewFfmpeg(log *logger.Logger) (*Ffmpeg, error) {
	return &Ffmpeg{log: log}, nil
}

type Ffmpeg struct {
	io.Writer
	log     *logger.Logger
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	started bool
	out     io.Writer
	buffer  *bytes.Buffer
}

func (f *Ffmpeg) SetOutput(out io.Writer) error {
	f.out = out
	f.cmd = exec.Command(
		"ffmpeg",
		"-y", // Yes to all
		"-hide_banner", "-loglevel", "panic", // Hide all logs
		"-i", "pipe:0", // take stdin as input
		// "-i", "pipe:", // take stdin as input
		"-map_metadata", "-1", // strip out all (mostly) metadata
		// "-c:a", "libopus", // use mp3 lame codec
		// "-vsync", "2", // suppress "Frame rate very high for a muxer not efficiently supporting it"
		// "-acodec", "pcm_s161e",
		// "-ar", "44100",
		// "-b:a", "128k", // Down sample audio birate to 128k
		"-ar", "16000", "-ac", "1", "-acodec", "pcm_s16le",
		"-f", "wav", // using mp3 muxer (IMPORTANT, output data to pipe require manual muxer selecting)
		"pipe:1", // output to stdout
		// "pipe:", // output to stdout
		// "/home/renan/projetos/doctor_recorder/teste.wav",
		// "output.mp3",
		// "|", "./audio.mp3",
	)
	var outb bytes.Buffer
	// f.cmd.Stdout = os.Stdout
	f.cmd.Stderr = os.Stderr
	f.cmd.Stdout = out
	f.buffer = &outb
	stdin, err := f.cmd.StdinPipe() // Open stdin pipe
	if err != nil {
		return err
	}
	// stdout, err := f.cmd.StdoutPipe() // Open stdin pipe
	// if err != nil {
	// 	return err
	// }
	// f.stdout = stdout
	f.stdin = stdin
	f.started = false
	return nil
}

func (f *Ffmpeg) Write(buf []byte) (int, error) {
	// var resultBuffer bytes.Buffer // pre allocate 5MiB buffer
	// var buff []byte
	// f.cmd.Stdout = &resultBuffer
	if f.started != true {
		f.log.Info("starting ffmpeg process")
		err := f.cmd.Start()
		if err != nil {
			f.log.Error(fmt.Errorf("failed to start ffmpeg command : %s", err).Error())
			return 0, err
		}
		f.started = true
	}
	// f.cmd.Stdout = &resultBuffer
	// f.log.Info(fmt.Sprintf("sending bytes %v", buf))
	// write to stding
	// io.WriteString(f.stdin, string(buf))
	f.stdin.Write(buf)
	// io.ReadFull(f.stdout, buf)
	// f.cmd.Wait()
	// f.log.Info(fmt.Sprintf("sending bytes %v", resultBuffer.Bytes()))
	// f.stdout.Read(buff)
	// f.log.Info(fmt.Sprintf("input bytes: %v", buf))
	// f.log.Info(fmt.Sprintf("output bytes: %v", resultBuffer))
	// f.out.Write(f.buffer.Bytes())
	// f.buffer.Reset()
	// time.Sleep(2 * time.Second)
	return 0, nil
}

func (f *Ffmpeg) Close() error {
	f.log.Warning("closing ffmpeg connection")
	f.stdin.Close()
	f.cmd.Wait()
	// err := f.stdin.Close() // close the stdin, or ffmpeg will wait forever
	// if err != nil {
	// 	return err
	// }
	// err = f.cmd.Wait() // wait until ffmpeg finish
	// if err != nil {
	// 	return err
	// }
	return nil
}
