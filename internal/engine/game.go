package engine

import (
	"math/rand"
	"time"
)

type Point struct {
	X int
	Y int
}
type Game struct {
	Width              int
	Height             int
	Food               []Point
	Snake              []Point
	Direction          Point
	GameOver           bool
	Score              int
	FoodNearby         bool
	FoodCount          int
	ShrinkingFruits    []Point
	ShrinkingFoodCount int
	DirectionQueue     []Point
}

func absI(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func NewGame(width, height, count, scount int) *Game {
	snakeHead := Point{X: width/2 - 3, Y: height / 2}
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
		ShrinkingFoodCount: scount,
		ShrinkingFruits:    []Point{{snakeHead.X + 7, snakeHead.Y + 5}},
		FoodCount:          count,
		Food:               []Point{{snakeHead.X + 7, snakeHead.Y}},
		FoodNearby:         true,
	}
	g.initialPlaceFoods()
	g.initialPlaceShrinkingFoods()
	return g
}
func (g *Game) initialPlaceFoods() {
	emptyCells := getEmptyCells(g)
	for i := 0; i < g.FoodCount-1; i++ {
		rand.NewSource(time.Now().UnixNano())
		randomIndex := rand.Intn(len(emptyCells))
		g.Food = append(g.Food, emptyCells[randomIndex])
		emptyCells = append(emptyCells[:randomIndex], emptyCells[randomIndex+1:]...)
	}
}
func (g *Game) initialPlaceShrinkingFoods() {
	emptyCells := getEmptyCells(g)
	for i := 0; i < g.ShrinkingFoodCount-1; i++ {
		rand.NewSource(time.Now().UnixNano())
		randomIndex := rand.Intn(len(emptyCells))
		g.ShrinkingFruits = append(g.ShrinkingFruits, emptyCells[randomIndex])
		emptyCells = append(emptyCells[:randomIndex], emptyCells[randomIndex+1:]...)
	}
}
func getEmptyCells(g *Game) []Point {
	width, height := g.Width, g.Height
	var emptyCells []Point
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			currentPoint := Point{X: x, Y: y}
			if !checkSnake(g, currentPoint) && !checkFood(g, currentPoint) && !checkShrinkingFruit(g, currentPoint) {
				emptyCells = append(emptyCells, currentPoint)
			}
		}
	}
	return emptyCells
}
func (g *Game) placeFood() {
	emptyCells := getEmptyCells(g)
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(emptyCells))
	g.Food = append(g.Food, emptyCells[randomIndex])
}
func (g *Game) placeShrinkingFruit() {
	emptyCells := getEmptyCells(g)
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(emptyCells))
	g.ShrinkingFruits = append(g.ShrinkingFruits, emptyCells[randomIndex])
}
func checkSnake(g *Game, point Point) bool {
	for _, snakeHead := range g.Snake {
		if snakeHead.X == point.X && snakeHead.Y == point.Y {
			return true
		}
	}
	return false
}
func checkShrinkingFruit(g *Game, point Point) bool {
	for _, sfruit := range g.ShrinkingFruits {
		if sfruit.X == point.X && sfruit.Y == point.Y {
			return true
		}
	}
	return false
}
func checkFood(g *Game, point Point) bool {
	for _, food := range g.Food {
		if food.X == point.X && food.Y == point.Y {
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

	if g.checkWallCollision(newHead) || g.selfCollision(newHead) {
		g.GameOver = true
		return
	}
	if g.checkFoodNearby(newHead) {
		g.FoodNearby = true
	} else {
		g.FoodNearby = false
	}
	g.Snake = append([]Point{newHead}, g.Snake...)
	eaten, check := eatenFood(g.Food, newHead)
	if check {
		g.Score++
		g.placeFood()
		g.Food = append(g.Food[:eaten], g.Food[eaten+1:]...)
	} else {
		g.Snake = g.Snake[:len(g.Snake)-1]
	}
	seaten, scheck := eatenFood(g.ShrinkingFruits, newHead)
	if scheck {
		g.Score--
		g.placeShrinkingFruit()
		g.ShrinkingFruits = append(g.ShrinkingFruits[:seaten], g.ShrinkingFruits[seaten+1:]...)
		g.Snake = g.Snake[:len(g.Snake)-1]
	}
	if len(g.Snake) == 0 {
		g.GameOver = true
	}

}
func eatenFood(food []Point, point Point) (int, bool) {
	for i, f := range food {
		if f.X == point.X && f.Y == point.Y {
			return i, true
		}
	}
	return -1, false
}
func (g *Game) UpdateDirectionQueue(x, y int) {
	directionx := g.Direction.X
	directiony := g.Direction.Y
	if directionx+x != 0 || directiony+y != 0 && g.Direction != g.Direction && len(g.DirectionQueue) < 3 {
		g.DirectionQueue = append(g.DirectionQueue, Point{X: x, Y: y})
	}
}
func (g *Game) checkWallCollision(p Point) bool {
	if p.X < 0 || p.Y < 0 || p.X >= g.Width || p.Y >= g.Height {
		return true
	}
	return false
}
func (g *Game) selfCollision(p Point) bool {
	for i := 1; i < len(g.Snake); i++ {
		if p.X == g.Snake[i].X && p.Y == g.Snake[i].Y {
			return true
		}
	}
	return false
}
func (g *Game) checkFoodNearby(p Point) bool {
	px, py := p.X, p.Y
	for _, fruit := range g.Food {
		gx, gy := fruit.X, fruit.Y
		if maxI(absI(gx-px), absI(gy-py)) <= 2 {
			return true
		}
	}
	return false
}
