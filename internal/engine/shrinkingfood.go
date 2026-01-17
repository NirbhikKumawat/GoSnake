package engine

import (
	"math/rand"
	"time"
)

func (g *Game) initialPlaceShrinkingFoods() {
	emptyCells := getEmptyCells(g)
	for i := 0; i < g.ShrinkingFoodCount-1; i++ {
		rand.NewSource(time.Now().UnixNano())
		randomIndex := rand.Intn(len(emptyCells))
		g.ShrinkingFruits = append(g.ShrinkingFruits, emptyCells[randomIndex])
		emptyCells = append(emptyCells[:randomIndex], emptyCells[randomIndex+1:]...)
	}
}
func (g *Game) placeShrinkingFruit() {
	emptyCells := getEmptyCells(g)
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(emptyCells))
	g.ShrinkingFruits = append(g.ShrinkingFruits, emptyCells[randomIndex])
}
func checkShrinkingFruit(g *Game, point Point) bool {
	for _, sfruit := range g.ShrinkingFruits {
		if sfruit.X == point.X && sfruit.Y == point.Y {
			return true
		}
	}
	return false
}
