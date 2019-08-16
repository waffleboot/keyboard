package main

import (
	"log"
	"os"
	"time"
	"unicode"

	"github.com/pkg/term"

	"github.com/faiface/beep/mp3"

	"github.com/faiface/beep/speaker"

	"github.com/faiface/beep"
)

func ReadSound(filename string) (beep.StreamSeekCloser, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	streamer, _, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}
	return streamer, nil
}

type Sounds map[rune]beep.StreamSeekCloser

func (s Sounds) Clear() {
	for _, v := range s {
		v.Close()
	}
}

func (s Sounds) Close() {
	s.Clear()
}

func ReadSounds() (Sounds, error) {
	sounds := Sounds(make(map[rune]beep.StreamSeekCloser))
	for c := 'A'; c <= 'Z'; c++ {
		streamer, err := ReadSound("sounds/" + string(c) + ".mp3")
		if err != nil {
			sounds.Clear()
			return nil, err
		}
		sounds[c] = streamer
	}
	return sounds, nil
}

func (s Sounds) Play(c rune) {
	done := make(chan struct{})
	streamer := s[c]
	streamer.Seek(0)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))
	<-done
}

func PlayAndWait(s Sounds, t *term.Term, c rune) {
	if c < 'A' || 'Z' < c {
		return
	}
	var buf [1]byte
	for c != unicode.ToUpper(rune(buf[0])) {
		s.Play(c)
		_, err := t.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	sampleRate := beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	sounds, err := ReadSounds()
	if err != nil {
		log.Fatal(err)
	}
	defer sounds.Close()
	text := "It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like)."
	t, err := term.Open("/dev/tty")
	if err != nil {
		log.Fatal(err)
	}
	defer t.Restore()
	t.SetCbreak()
	for _, c := range text {
		c = unicode.ToUpper(c)
		PlayAndWait(sounds, t, c)
	}
}
