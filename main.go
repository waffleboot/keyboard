package main

import (
	"fmt"
	"log"
	"math/rand"
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
	var retry int
	var buf [1]byte
	for c != unicode.ToUpper(rune(buf[0])) {
		if retry > 2 {
			fmt.Printf("%q\n", c)
			retry = 0
		}
		s.Play(c)
		_, err := t.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		retry++
	}
}

func main() {
	rand.Seed(42)
	sampleRate := beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	sounds, err := ReadSounds()
	if err != nil {
		log.Fatal(err)
	}
	defer sounds.Close()
	t, err := term.Open("/dev/tty")
	if err != nil {
		log.Fatal(err)
	}
	defer t.Restore()
	t.SetCbreak()
	for {
		c := unicode.ToUpper(rune('A' + rand.Intn(26)))
		PlayAndWait(sounds, t, c)
	}
}