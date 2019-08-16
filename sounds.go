package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/faiface/beep/speaker"

	"github.com/faiface/beep/mp3"

	"github.com/faiface/beep"
)

type Sounds map[rune]beep.StreamSeekCloser

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

func (s Sounds) Clear() {
	for _, v := range s {
		v.Close()
	}
}

func (s Sounds) Close() {
	s.Clear()
}
