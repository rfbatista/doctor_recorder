package ffmpeg

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

func NewFfmpeg() (*Ffmpeg, error) {
	return &Ffmpeg{}, nil
}

type Ffmpeg struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

func (f *Ffmpeg) Init() {
	f.cmd = exec.Command("ffmpeg", "-y", // Yes to all
		//"-hide_banner", "-loglevel", "panic", // Hide all logs
		"-i", "pipe:0", // take stdin as input
		"-map_metadata", "-1", // strip out all (mostly) metadata
		"-c:a", "libmp3lame", // use mp3 lame codec
		"-vsync", "2", // suppress "Frame rate very high for a muxer not efficiently supporting it"
		"-b:a", "128k", // Down sample audio birate to 128k
		"-f", "mp3", // using mp3 muxer (IMPORTANT, output data to pipe require manual muxer selecting)
		"pipe:1", // output to stdout
	)
}

func (f *Ffmpeg) Read(buf []byte) (*bytes.Buffer, error) {
	resultBuffer := bytes.NewBuffer(make([]byte, 5*1024*1024)) // pre allocate 5MiB buffer
	f.cmd.Stdout = os.Stdout
	f.cmd.Stderr = os.Stderr
	stdin, err := f.cmd.StdinPipe() // Open stdin pipe
	f.stdin = stdin
	if err != nil {
		return nil, err
	}
	err = f.cmd.Start()
	if err != nil {
		return nil, err
	}
	_, err = stdin.Write(buf) // pump audio data to stdin pipe
	if err != nil {
		return nil, err
	}
	return resultBuffer, nil
}

func (f *Ffmpeg) Close() error {
	err := f.stdin.Close() // close the stdin, or ffmpeg will wait forever
	if err != nil {
		return err
	}
	err = f.cmd.Wait() // wait until ffmpeg finish
	if err != nil {
		return err
	}
	return nil
}
