package main

import (
	"math/rand"
	"time"
)

type Point struct {
	X int
	Y int
}
type Game struct {
	Width     int
	Height    int
	Food      Point
	Snake     []Point
	Direction Point
	GameOver  bool
	Score     int
}

func NewGame(width, height int) *Game {
	snakeHead := Point{X: width / 2, Y: height / 2}
	g := &Game{
		Width:    width,
		Height:   height,
		GameOver: false,
		Score:    0,
		Direction: Point{
			X: 1,
			Y: 0,
		},
		Snake: []Point{
			snakeHead,
			{
				X: snakeHead.X - 1,
				Y: snakeHead.Y,
			},
			{
				X: snakeHead.X - 2,
				Y: snakeHead.Y,
			},
		},
		Food: Point{snakeHead.X + 2, snakeHead.Y},
	}
	return g
}
func (g *Game) PlaceFood() {
	width, height := g.Width, g.Height
	var emptyCells []Point
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			currentPoint := Point{X: x, Y: y}
			if !CheckSnake(g, currentPoint) {
				emptyCells = append(emptyCells, currentPoint)
			}
		}
	}
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(emptyCells))
	g.Food = emptyCells[randomIndex]
}
func CheckSnake(g *Game, point Point) bool {
	for _, snakeHead := range g.Snake {
		if snakeHead.X == point.X && snakeHead.Y == point.Y {
			return true
		}
	}
	return false
}
func (g *Game) Move() {
	if g.GameOver {
		return
	}
	newHead := Point{X: g.Snake[0].X + g.Direction.X, Y: g.Snake[0].Y + g.Direction.Y}

	if g.CheckWallCollision(newHead) || g.SelfCollision(newHead) {
		g.GameOver = true
		return
	}
	g.Snake = append([]Point{newHead}, g.Snake...)
	if g.Food.X == newHead.X && g.Food.Y == newHead.Y {
		g.Score++
		g.PlaceFood()
	} else {
		g.Snake = g.Snake[:len(g.Snake)-1]
	}
}
func (g *Game) NextDirection(x, y int) Point {
	directionx := g.Direction.X
	directiony := g.Direction.Y
	if directionx+x == 0 && directiony+y == 0 {
		return g.Direction
	}
	return Point{X: x, Y: y}
}
func (g *Game) CheckWallCollision(p Point) bool {
	if p.X < 0 || p.Y < 0 || p.X >= g.Width || p.Y >= g.Height {
		return true
	}
	return false
}
func (g *Game) SelfCollision(p Point) bool {
	for i := 1; i < len(g.Snake); i++ {
		if p.X == g.Snake[i].X && p.Y == g.Snake[i].Y {
			return true
		}
	}
	return false
}
