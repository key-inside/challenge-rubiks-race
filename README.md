# Rubik's Race
![real playing](./images/rubiks_race.gif)
- [About](https://boardgamegeek.com/wiki/page/thing:614)
- [Online game](https://sites.google.com/view/azunblockedgames/rubiks-race)

# How to challange

## Be a gopher
- [The Go Gopher](https://go.dev/blog/gopher)

## Write racer plugin
1. Make a plugin source code file in `./plugins/{your_name}/` (ex, ./plugins/[max](https://en.wikipedia.org/wiki/Max_Park)/main.go)
2. Implement interfaces below
```go
import (
    "rubiks-race/pkg/puzzle"
    "rubiks-race/pkg/racer"
)

func NewRacer() racer.Racer
```

## Build racer plugin
```sh
% PLUGIN_NAME=max make build-plugin
```

## Build game
```sh
% make build
```

## Run game
### Run with random puzzle
```sh
% PLUGIN_NAME=max make run
# or
% ./build/rubiks-race max
```
### Run with custom puzzle
> PUZZLE_DATA={24bytes-puzzle}:{9bytes-goal}
```sh
% PLUGIN_NAME=max PUZZLE_DATA="GROWBYGROWBYGROWBYGROWBY:GROWBYRGB" make run
# or
% ./build/rubiks-race max "GROWBYGROWBYGROWBYGROWBY:GROWBYRGB"
```

## Follow logs
```sh
% tail -f rubiks.log
```

## What baby gophers should study
- https://github.com/golang-standards/project-layout/blob/master/README_ko.md
- https://golang.design/research/generic-option/
- https://yourbasic.org/golang/compare-slices/
