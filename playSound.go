package main

import (
	"bytes"
	_ "embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"log"
	"time"
)


//go:embed suffer.mp3
var audio []byte

type Player struct {
	AudioBuffer *beep.Buffer
	Format beep.Format
}


func (p *Player) PlayAudio()  {
	done := make(chan bool, 1)
	streamer := p.AudioBuffer.Streamer(0, p.AudioBuffer.Len())
	speaker.Play(beep.Seq(streamer), beep.Callback(func() {
		done <- true
	}))
	<- done
}

func createAudioBuffer() (*beep.Buffer, beep.Format) {
	r := io.NopCloser(bytes.NewReader(audio))
	streamer, format, err := mp3.Decode(r)
	if err != nil {
		log.Fatal(err)
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	return buffer, format
}


func NewPlayer() *Player {
	buffer, format := createAudioBuffer()
	err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))

	if err != nil {
		log.Fatal(err)
	}
	return &Player{ AudioBuffer: buffer, Format: format }
}
