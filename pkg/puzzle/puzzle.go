package puzzle

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type Cube byte

const (
	G = 'G'
	R = 'R'
	O = 'O'
	W = 'W'
	B = 'B'
	Y = 'Y'
	X = 'X'
)

func (c Cube) String() string {
	return string(c)
}

// var _cubes = []byte("GROWBYGROWBYGROWBYGROWBY")
var _cubes = []Cube{
	G, R, O, W, B, Y,
	G, R, O, W, B, Y,
	G, R, O, W, B, Y,
	G, R, O, W, B, Y,
}

type Direction int

const (
	Up = iota
	Right
	Down
	Left
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "UP"
	case Right:
		return "RIGHT"
	case Down:
		return "DOWN"
	case Left:
		return "LEFT"
	}
	return ""
}

type Option func(*Puzzle)

type Puzzle struct {
	cubes  []Cube
	goal   []Cube
	xIndex int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func scrambleCubes(size int) []Cube {
	rand.Shuffle(len(_cubes), func(i, j int) {
		_cubes[i], _cubes[j] = _cubes[j], _cubes[i]
	})
	cubes := make([]Cube, size)
	copy(cubes, _cubes)
	return cubes
}

func New(opts ...Option) *Puzzle {
	puzzle := &Puzzle{
		xIndex: 24,
	}

	for _, o := range opts {
		o(puzzle)
	}

	if nil == puzzle.goal { // creates random goal
		puzzle.goal = scrambleCubes(9)
	}

	if nil == puzzle.cubes { // scrambles puzzle
		puzzle.cubes = scrambleCubes(25)
		puzzle.cubes[24] = X
	}

	log.Printf("PUZZLE: cubes = %+v, goal = %+v", string(puzzle.cubes), string(puzzle.goal))

	return puzzle
}

func WithCubes(cubes []Cube) Option {
	// checks only length and does not validate values
	if len(cubes) != 25 {
		cubes = nil
	}
	return func(p *Puzzle) {
		p.cubes = cubes
	}
}

func WithGoal(goal []Cube) Option {
	// checks only length and does not validate values
	if len(goal) != 9 {
		goal = nil
	}
	return func(p *Puzzle) {
		p.goal = goal
	}
}

func (p *Puzzle) IsComplete() bool {
	return (string(p.cubes[6:9]) == string(p.goal[0:3])) &&
		(string(p.cubes[11:14]) == string(p.goal[3:6])) &&
		(string(p.cubes[16:19]) == string(p.goal[6:9]))
}

func (p *Puzzle) Cubes() []Cube {
	return p.cubes
}

func (p *Puzzle) Goal() []Cube {
	return p.goal
}

func (p *Puzzle) Move(i int) (Direction, error) {
	var d Direction
	if i > 0 && i-1 == p.xIndex {
		d = Left
	} else if i > 4 && i-5 == p.xIndex {
		d = Up
	} else if i < 24 && i+1 == p.xIndex {
		d = Right
	} else if i < 20 && i+5 == p.xIndex {
		d = Down
	} else {
		return d, fmt.Errorf("the cube at %d can't move", i)
	}

	p.cubes[i], p.cubes[p.xIndex] = p.cubes[p.xIndex], p.cubes[i]
	p.xIndex = i
	return d, nil
}

func (p *Puzzle) CubeAt(i int) Cube {
	return p.cubes[i]
}

func (p *Puzzle) GoalAt(i int) Cube {
	return p.goal[i]
}

func (p *Puzzle) IndexOfX() int {
	return p.xIndex
}
