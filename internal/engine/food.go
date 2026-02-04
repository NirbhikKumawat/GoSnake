package engine

import (
	"math/rand"
	"time"
)

func (g *Game) initialPlaceFoods() {
	emptyCells := getEmptyCells(g)
	for i := 0; i < g.FoodCount-1; i++ {
		rand.NewSource(time.Now().UnixNano())
		randomIndex := rand.Intn(len(emptyCells))
		g.Food = append(g.Food, emptyCells[randomIndex])
		emptyCells = append(emptyCells[:randomIndex], emptyCells[randomIndex+1:]...)
	}
}
func (g *Game) placeFood() {
	emptyCells := getEmptyCells(g)
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(emptyCells))
	g.Food = append(g.Food, emptyCells[randomIndex])
}
func checkFood(g *Game, point Point) bool {
	for _, food := range g.Food {
		if food.X == point.X && food.Y == point.Y {
			return true
		}
	}
	return false
}
func eatenFood(food []Point, point Point) (int, bool) {
	for i, f := range food {
		if f.X == point.X && f.Y == point.Y {
			return i, true
		}
	}
	return -1, false
}
func (g *Game) handleEatenFood(eaten int) {
	g.Score++
	g.placeFood()
	if g.Score%2 == 0 && g.SpawnBlocks {
		g.placeBlocks()
	}
	if g.ReverseDirection {
		g.Direction = g.getDirection()
		for i, j := 0, len(g.Snake)-1; i < j; i, j = i+1, j-1 {
			g.Snake[i], g.Snake[j] = g.Snake[j], g.Snake[i]
		}
	}
	g.Food = append(g.Food[:eaten], g.Food[eaten+1:]...)
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
