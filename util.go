package main

import "math/rand"

type LettersGenerator rune

func init() {
	rand.Seed(42)
}

func (g *LettersGenerator) Next() rune {
	c := rune('a' + rand.Intn(26))
	for c == rune(*g) {
		c = rune('a' + rand.Intn(26))
	}
	*g = LettersGenerator(c)
	return c
}
