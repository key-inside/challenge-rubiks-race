package game

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"

	"rubiks-race/pkg/puzzle"
	"rubiks-race/pkg/racer"
)

const (
	fps     = 10                // frame per second
	fiv     = time.Second / fps // frame interval
	fpa     = 1                 // frame per animation
	fpm     = fpa * 3           // frame per move
	timeout = time.Second * 5
)

type Status int

const (
	Running = iota
	Complete
	Retire
	Timeout
)

type Option func(*Game)

type Game struct {
	sync.Mutex

	screen    tcell.Screen
	startTime time.Time
	frame     int
	try       int
	status    Status
	isStop    bool

	puzzle    *puzzle.Puzzle
	direction puzzle.Direction
	cIndex    int

	racer racer.Racer
}

func New(racr racer.Racer, pzl *puzzle.Puzzle, opts ...Option) (*Game, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err

	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	// default style
	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))

	g := &Game{
		screen: s,
		puzzle: pzl,
		racer:  racr,
	}

	return g, nil
}

func (g *Game) resizeScreen() {
	g.Lock()
	g.screen.Sync()
	g.Unlock()
}

func (g *Game) Quit() {
	g.Lock()
	g.screen.Fini()
	g.Unlock()
	os.Exit(0)
}

func (g *Game) Run() {
	defer g.Quit()

	go g.listenUserEvents()

	g.startTime = time.Now()

	ticker := time.NewTicker(fiv)
	defer ticker.Stop()

	go g.race()

	for range ticker.C {
		if !g.isStop {
			if g.status != Running {
				g.isStop = true
			}
			g.updateScreen()
		}
	}
}

func (g *Game) race() {
	for g.status == Running {
		for g.frame < fpm { // waiting cube animation
			time.Sleep(fiv)
		}

		if g.puzzle.IsComplete() {
			g.status = Complete
			log.Printf("COMPLETE: Try Count = %d, Elapsed Time = %s", g.try, time.Since(g.startTime).String())
		} else {
			res := make(chan int, 1)
			go func() {
				if g.try == 0 {
					g.racer.Init(g.puzzle.Cubes(), g.puzzle.Goal())
				}
				g.try++
				res <- g.racer.Next(g.try)
			}()

			select {
			case ci := <-res:
				// replace current moving cube index with current X index
				g.cIndex = g.puzzle.IndexOfX()
				var err error
				g.direction, err = g.puzzle.Move(ci)
				if err != nil {
					log.Println("ERROR:", err.Error())
					g.status = Retire
					break
				}

				log.Printf("MOVE: Try Count = %d, Cube Color/Index = %s/%d, Direction = %s",
					g.try,
					g.puzzle.CubeAt(g.cIndex),
					ci,
					g.direction)

				g.frame = 0 // reset animation frame
			case <-time.After(timeout):
				log.Printf("ERROR: timeout! move next cube in %s", timeout.String())
				g.status = Timeout
			}
		}
	}
}

func (g *Game) listenUserEvents() {
	for {
		ev := g.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			g.resizeScreen()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				g.Quit()
			}
		}
	}
}

func (g *Game) drawText(x1, y1, x2, y2 int, text string, style tcell.Style) {
	row := y1
	col := x1
	for _, r := range text {
		g.screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (g *Game) updateScreen() {
	g.screen.Clear()
	g.drawHelp()
	g.drawMetric()
	g.drawMsg()
	g.drawBoard()
	g.screen.Show()
}

var msgStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
var completeMsgStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
var errorMsgStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
var metricNameStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
var metricValueStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

func (g *Game) drawHelp() {
	g.drawText(5, 1, 34, 1, "(Press ESC or Ctrl+C to quit)", msgStyle)
}

func (g *Game) drawMetric() {
	g.drawText(5, 3, 18, 3, "Elapsed Time:", metricNameStyle)
	g.drawText(19, 3, 44, 3, time.Since(g.startTime).String(), metricValueStyle)
	g.drawText(5, 4, 15, 4, "Try Count:", metricNameStyle)
	g.drawText(16, 4, 26, 4, strconv.Itoa(g.try), metricValueStyle)
}

func (g *Game) drawMsg() {
	switch g.status {
	case Complete:
		g.drawText(15, 6, 33, 6, "    Completed!    ", completeMsgStyle)
	case Retire:
		g.drawText(15, 6, 33, 6, "     Retire!      ", errorMsgStyle)
	case Timeout:
		g.drawText(15, 6, 33, 6, "     Timeout!     ", errorMsgStyle)
	}
}

func colorOf(cube puzzle.Cube) tcell.Color {
	switch cube {
	case puzzle.G:
		return tcell.ColorGreen
	case puzzle.R:
		return tcell.ColorRed
	case puzzle.O:
		return tcell.ColorOrange
	case puzzle.W:
		return tcell.ColorWhite
	case puzzle.B:
		return tcell.ColorBlue
	case puzzle.Y:
		return tcell.ColorYellow
	}
	return tcell.ColorBlack
}

func (g *Game) drawBoard() {
	g.drawCubes()
	g.drawGoal()
}

func (g *Game) drawCubes() {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			ci := i*5 + j
			cube := g.puzzle.CubeAt(ci)
			if puzzle.X == cube {
				continue
			}

			style := tcell.StyleDefault.Background(colorOf(cube))
			x := 5 + j*8
			y := 8 + i*4

			// animating
			if g.status == Running {
				if ci == g.cIndex && g.frame < fpm {
					aniFrame := (g.frame % fpm) / fpa
					d := 2 - aniFrame
					if g.try > 0 {
						switch g.direction {
						case puzzle.Up:
							y += d
						case puzzle.Right:
							x -= d * 2
						case puzzle.Down:
							y -= d
						case puzzle.Left:
							x += d * 2
						}
					}
					g.frame++
				}
			}

			g.drawText(x, y, x+6, y+2, "                  ", style)
		}
	}
}

func (g *Game) drawGoal() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			gi := i*3 + j
			style := tcell.StyleDefault.Background(colorOf(g.puzzle.GoalAt(gi)))
			x := 15 + j*8
			y := 13 + i*4
			g.drawText(x, y, x+2, y, "  ", style)
		}
	}
}
