package main

import (
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
	quit          chan struct{}
}

func NewGame(width, height int) *Game {
	return &Game{
		snake:  []Point{{20, 8}},
		dir:    Point{1, 0},
		width:  width,
		height: height,
		quit:   make(chan struct{}),
	}
}

func (g *Game) draw(ev tb.Event) {
	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	tb.Flush()
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

	snakeHead := g.snake[0]
	tb.SetCell(snakeHead.x, snakeHead.y, g.dir.ToRune(ev), tb.ColorRed, tb.ColorDefault)

	for i, ch := range "Score: 0" {
		tb.SetCell(i+1, 0, ch, tb.ColorBlack, tb.ColorDefault)
	}

	tb.Flush()

}

func (g *Game) handleInput(ev tb.Event) {
	key := g.dir.ToRune(ev)

	switch key {
	case '↑':
		g.dir.y += 1
	case '↓':
		g.dir.y -= 1
	case '←':
		g.dir.x -= 1
	case '→':
		g.dir.x += 1
	default:
		tb.Close()
	}
}

func (p *Point) ToRune(ev tb.Event) rune {
	key := ev.Key

	switch true {
	case key == tb.KeyArrowUp || key == tb.KeyCtrlW:
		return '↑'
	case key == tb.KeyArrowDown || key == tb.KeyCtrlS:
		return '↓'
	case key == tb.KeyArrowLeft || key == tb.KeyCtrlA:
		return '←'
	default:
		return '→'
	}
}

func main() {
	tb.Init()
	defer tb.Close()

	eventQueue := make(chan tb.Event)

	go func() {
		for {
			eventQueue <- tb.PollEvent()
		}
	}()

	game := NewGame(25, 18)
	for {
		ev := <-eventQueue
		game.draw(ev)
		game.handleInput(ev)
	}
}
