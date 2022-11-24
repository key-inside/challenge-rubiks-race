package main

import (
	"log"
	"math/rand"
	"rubiks-race/pkg/puzzle"
	"rubiks-race/pkg/racer"
	"time"
)

type _racer struct {
	x int
}

func NewRacer() racer.Racer {
	return &_racer{}
}

func (r *_racer) Init(cubes []puzzle.Cube, goal []puzzle.Cube) {
	for i, cube := range cubes {
		if puzzle.X == cube {
			r.x = i
			break
		}
	}
	log.Print("Racer Rooney: ", "Initialzed, X =", r.x)
}

func (r *_racer) Next(try int) int {
	rand.Seed(time.Now().UnixNano())

	if try%2 == 0 {
		if r.x < 5 {
			r.x += 5 // move up
		} else if r.x > 19 {
			r.x -= 5 // move down
		} else if rand.Int()%2 == 0 {
			r.x += 5 // move up
		} else {
			r.x -= 5 // move down
		}
	} else {
		col := r.x % 5
		if col == 0 {
			r.x += 1 // move left
		} else if col == 4 {
			r.x -= 1 // move right
		} else if rand.Int()%2 == 0 {
			r.x += 1 // move left
		} else {
			r.x -= 1 // move right
		}
	}

	log.Print("Racer Rooney: ", "Try = ", try, ", Next = ", r.x)
	return r.x
}
