package engine

import (
	"math/rand"
	"time"
)

func (g *Game) placePortalFruits() {
	emptyCells := getEmptyCells(g)
	rand.NewSource(time.Now().UnixNano())
	randomIndex1 := rand.Intn(len(emptyCells))
	cell1 := emptyCells[randomIndex1]
	emptyCells = append(emptyCells[:randomIndex1], emptyCells[randomIndex1+1:]...)
	randomIndex2 := rand.Intn(len(emptyCells))
	cell2 := emptyCells[randomIndex2]
	g.PortalFruits = append(g.PortalFruits, [2]Point{cell1, cell2})
}
func checkPortalFoods(g *Game, point Point) bool {
	for _, food := range g.PortalFruits {
		if food[0].X == point.X && food[0].Y == point.Y {
			return true
		}
		if food[1].X == point.X && food[1].Y == point.Y {
			return true
		}
	}
	return false
}
func eatenPortalFood(foods [][2]Point, point Point) (int, int, bool) {
	for i, food := range foods {
		if food[0].X == point.X && food[0].Y == point.Y {
			return 1, i, true
		}
		if food[1].X == point.X && food[1].Y == point.Y {
			return 0, i, true
		}
	}
	return -1, -1, false
}
