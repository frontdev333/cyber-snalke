package main

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"time"

	tb "github.com/nsf/termbox-go"
)

type Point struct {
	x int
	y int
}

type Game struct {
	snake         []Point
	food          Point
	malware       []Point
	dir           Point
	score         int
	gameOver      bool
	width, height int
	level         int
}

func NewGame(width, height int) *Game {
	g := &Game{
		snake:  []Point{{width/2 - 1, height/2 - 1}},
		dir:    Point{1, 0},
		width:  width,
		height: height,
		level:  1,
	}

	g.placeFood()
	g.placeMalware()

	return g
}

func (g *Game) draw() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(0, 0, '┌', tb.ColorBlack, tb.ColorDefault)
	tb.SetCell(g.width, 0, '┐', tb.ColorBlack, tb.ColorDefault)
	tb.SetCell(0, g.height, '└', tb.ColorBlack, tb.ColorDefault)
	tb.SetCell(g.width, g.height, '┘', tb.ColorBlack, tb.ColorDefault)

	for i := 1; i < g.width; i++ {
		tb.SetCell(i, 0, '─', tb.ColorBlack, tb.ColorDefault)
		tb.SetCell(i, g.height, '─', tb.ColorBlack, tb.ColorDefault)
	}

	for i := 1; i < g.height; i++ {
		tb.SetCell(0, i, '│', tb.ColorBlack, tb.ColorDefault)
		tb.SetCell(g.width, i, '│', tb.ColorBlack, tb.ColorDefault)
	}

	for k, v := range g.snake {
		if k == 0 {
			tb.SetCell(v.x, v.y, 'ဏ', tb.ColorRed, tb.ColorDefault)
			continue
		}
		tb.SetCell(v.x, v.y, 'ဳ', tb.ColorRed, tb.ColorDefault)
	}

	tb.SetCell(g.food.x, g.food.y, '●', tb.ColorBlue, tb.ColorDefault)
	for _, mPoint := range g.malware {
		tb.SetCell(mPoint.x, mPoint.y, '✗', tb.ColorRed, tb.ColorDefault)
	}

	drawText(fmt.Sprintf("Score: %d Level: %d", g.score, g.level), 0, 0, tb.ColorBlack)

	tb.Flush()

}

func (g *Game) handleInput(ev tb.Event) {
	key := ev.Key

	switch true {
	case key == tb.KeyArrowUp || key == tb.KeyCtrlW:
		if g.dir.y != 1 {
			g.dir = Point{0, -1}
		}
	case key == tb.KeyArrowDown || key == tb.KeyCtrlS:
		if g.dir.y != -1 {
			g.dir = Point{0, 1}
		}
	case key == tb.KeyArrowRight || key == tb.KeyCtrlD:
		if g.dir.x != -1 {
			g.dir = Point{1, 0}
		}
	case key == tb.KeyArrowLeft || key == tb.KeyCtrlA:
		if g.dir.x != 1 {
			g.dir = Point{-1, 0}
		}
	}
}

func (g *Game) isOnSnake(p Point) bool {
	for _, bodyPart := range g.snake {
		if bodyPart == p {
			return true
		}
	}
	return false
}

func (g *Game) isOnMalware(p Point) bool {
	for _, malware := range g.malware {
		if malware == p {
			return true
		}
	}
	return false
}

func (g *Game) isOutOfBounds(p Point) bool {
	if g.height <= p.y || p.y <= 0 || g.width <= p.x || p.x <= 0 {
		return true
	}
	return false
}

func (g *Game) placeFood() {
	for i := 0; i < 10; i++ {
		place := getRandPoint(g)

		if isPointFree(g, place) && place != g.food {
			g.food = place
			return
		}
	}
	slog.Warn("there's no free places for food")
}

func (g *Game) placeMalware() {
	for i := 0; i < 10; i++ {
		place := getRandPoint(g)

		if isPointFree(g, place) && place != g.food {
			g.malware = append(g.malware, place)
			return
		}
	}
	slog.Warn("there's no free places for malware")
}

func (g *Game) move() {
	g.level = g.score/5 + 1

	switch true {
	case g.level == 2 && len(g.malware) < 2:
		g.placeMalware()
	case g.level == 3 && len(g.malware) < 3:
		g.placeMalware()
	}

	snakeHead := Point{
		x: g.dir.x + g.snake[0].x,
		y: g.dir.y + g.snake[0].y,
	}

	if !isPointFree(g, snakeHead) {
		g.gameOver = true
		return
	}

	g.snake = append([]Point{snakeHead}, g.snake...)

	if snakeHead == g.food {
		g.score++
		g.placeFood()
		return
	}

	g.snake = g.snake[:len(g.snake)-1]

}

func (g *Game) drawGameOver() {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	centerX := g.width / 2
	centerY := g.height / 2

	gmeOvr := "GAME OVER"
	gmeOvrX := getCenterTextCoordinates(gmeOvr, centerX)

	yrScr := fmt.Sprintf("Your score: %d", g.score)
	yrScrX := getCenterTextCoordinates(yrScr, centerX)

	yrLvl := fmt.Sprintf("Your level: %d", g.level)
	yrLvlX := getCenterTextCoordinates(yrLvl, centerX)

	instrct := "Press R to restart"
	instrctX := getCenterTextCoordinates(instrct, centerX)

	drawText(gmeOvr, gmeOvrX, centerY-2, tb.ColorRed)
	drawText(yrScr, yrScrX, centerY, tb.ColorBlack)
	drawText(yrLvl, yrLvlX, centerY+2, tb.ColorBlack)
	drawText(instrct, instrctX, centerY+4, tb.ColorGreen)

	tb.Flush()
}

func getCenterTextCoordinates(text string, centerX int) int {
	startX := centerX - (len(text) / 2)
	return startX
}

func (p *Point) ToRune() rune {
	switch true {
	case p.x == 0 && p.y == -1:
		return '↑'
	case p.x == 0 && p.y == 1:
		return '↓'
	case p.x == -1 && p.y == 0:
		return '←'
	default:
		return '→'
	}
}

func isPointFree(g *Game, place Point) bool {
	return !g.isOnSnake(place) && !g.isOnMalware(place) && !g.isOutOfBounds(place)
}

func getRandPoint(g *Game) Point {
	x := rand.IntN(g.width-2) + 1
	y := rand.IntN(g.height-2) + 1

	place := Point{x, y}
	return place
}

func drawText(text string, width, height int, color tb.Attribute) {
	for i, ch := range text {
		tb.SetCell(i+width+1, height, ch, color, tb.ColorDefault)
	}
}

func handleGameOver(game *Game, eventQueue chan tb.Event) bool {
	game.drawGameOver()
	for {
		ev := <-eventQueue

		if ev.Key == tb.KeyEsc || ev.Key == tb.KeyCtrlQ {
			return true
		}
		if ev.Ch == 'R' || ev.Ch == 'r' {
			*game = *NewGame(40, 18)
			break
		}
	}
	return false
}

func playGame() {
	err := tb.Init()
	if err != nil {
		fmt.Println("Ошибка инициализации termbox:", err)
		return
	}
	defer tb.Close()

	game := NewGame(40, 18)

	baseInterval := 200 * time.Millisecond
	speedModifier := 1.0 - float64(game.level-1)*0.1
	interval := time.Duration(float64(baseInterval) * max(speedModifier, 0.2))
	ticker := time.NewTicker(baseInterval)
	ticker.Reset(interval)

	defer ticker.Stop()

	eventQueue := make(chan tb.Event)

	go func() {
		for {
			eventQueue <- tb.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Key == tb.KeyEsc || ev.Key == tb.KeyCtrlQ {
				if handleGameOver(game, eventQueue) {
					return
				}
			}
			game.handleInput(ev)
		case <-ticker.C:
			if game.gameOver {
				handleGameOver(game, eventQueue)
			}
			if !game.gameOver {
				game.move()
			}
			game.draw()
		}
	}

}

func main() {
	playGame()
}
