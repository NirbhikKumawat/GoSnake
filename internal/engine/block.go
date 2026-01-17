package engine

import (
	"math/rand"
	"time"
)

func (g *Game) placeBlocks() {
	emptyCells := getEmptyCells(g)
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(emptyCells))
	g.Blocks = append(g.Blocks, emptyCells[randomIndex])
}
func checkBlocks(g *Game, point Point) bool {
	for _, block := range g.Blocks {
		if block.X == point.X && block.Y == point.Y {
			return true
		}
	}
	return false
}
