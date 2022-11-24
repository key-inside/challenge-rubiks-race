package main

import (
	"log"
	"os"
	"path"
	"plugin"
	"strings"

	"rubiks-race/pkg/game"
	"rubiks-race/pkg/puzzle"
	"rubiks-race/pkg/racer"
)

// args[1] = racer plugin name
// args[2] ?= "{24bytes-puzzle}:{9bytes-goal}" => "GROWBYGROWBYGROWBYGROWBY:GROWBYRGB"
func main() {
	if len(os.Args) < 2 {
		log.Fatalln("no racer plugin name!")
	}
	plgin, err := plugin.Open(path.Join("build/plugins/", os.Args[1]+".so"))
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	newFnSym, err := plgin.Lookup("NewRacer")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	newRacerFn, ok := newFnSym.(func() racer.Racer)
	if !ok {
		log.Fatalf("plugin does not implement function NewRacer")
	}

	// file output logging
	fpLog, err := os.OpenFile("rubiks.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%+v\n", err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)

	// creates puzzle
	opts := []puzzle.Option{}
	if len(os.Args) > 2 {
		data := strings.Split(os.Args[2], ":")
		if len(data) < 2 {
			log.Fatalln("wrong input puzzle data")
		}
		opts = append(opts,
			puzzle.WithCubes(append([]puzzle.Cube(data[0]), puzzle.X)),
			puzzle.WithGoal([]puzzle.Cube(data[1])))
	}

	// creates & starts game
	g, err := game.New(newRacerFn(), puzzle.New(opts...))
	if err != nil {
		log.Fatalf("%+v", err)
	}
	g.Run()
}
