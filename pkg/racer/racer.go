package racer

import "rubiks-race/pkg/puzzle"

type Racer interface {
	Init([]puzzle.Cube, []puzzle.Cube)
	Next(int) int
}
