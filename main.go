package main

import (
	"bytes"
	"fmt"
	"log"

	"time"

	"github.com/pkg/term"

	"github.com/faiface/beep/speaker"

	"github.com/faiface/beep"
)

func PlayAndWait(s Sounds, t *term.Term, c rune) {
	var retry int
	var buf [4]byte
	for c != bytes.Runes(buf[:])[0] {
		if retry > 1 {
			fmt.Printf("%q\n", c)
			retry = 0
		}
		if !s.Play(c) {
			return
		}
		_, err := t.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		retry++
	}
}

const sampleRate = beep.SampleRate(44100)

func init() {
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
}

func main() {
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
	var g LettersGenerator
	for {
		PlayAndWait(sounds, t, g.Next())
	}
}
