package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	// "sync/atomic"
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
	for c := 'a'; c <= 'z'; c++ {
		streamer, err := ReadSound(fmt.Sprintf("sounds/%s.mp3", string(unicode.ToUpper(c))))
		if err != nil {
			sounds.Clear()
			return nil, err
		}
		sounds[c] = streamer
	}
	return sounds, nil
}

func (s Sounds) Play(c rune) bool {
	streamer := s[c]
	if streamer == nil {
		return false
	}
	streamer.Seek(0)
	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))
	<-done
	return true
}

func PlayAndWait(s Sounds, t *term.Term, c rune) {
	var retry int
	var buf [1]byte
	for c != rune(buf[0]) {
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
	rand.Seed(42)
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
	g := &LettersGenerator{}
	for {
		PlayAndWait(sounds, t, g.Next())
	}
}

type LettersGenerator struct {
	prev rune
}

func (g *LettersGenerator) Next() rune {
	c := rune('a' + rand.Intn(26))
	for c == g.prev {
		c = rune('a' + rand.Intn(26))
	}
	g.prev = c
	return c
}

// func main() {
// 	rand.Seed(42)
// 	sounds, err := ReadSounds()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer sounds.Close()
// 	var letter rune
// 	g := &LettersGenerator{}
// 	go func() {
// 		for {
// 			c := g.Next()
// 			sounds.Play(c)
// 			for c != atomic.LoadInt32(&letter) {
// 				time.Sleep(time.Duration(250) * time.Millisecond)
// 				if c != atomic.LoadInt32(&letter) {
// 					sounds.Play(c)
// 				}
// 			}
// 		}
// 	}()
// 	go func() {
// 		t, err := term.Open("/dev/tty")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer t.Restore()
// 		t.SetCbreak()
// 		for {
// 			var buf [1]byte
// 			_, err := t.Read(buf[:])
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			atomic.StoreInt32(&letter, rune(buf[0]))
// 		}
// 	}()
// 	select {}
// }
