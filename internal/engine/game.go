package engine

type Point struct {
	X int
	Y int
}
type Game struct {
	Width               int
	Height              int
	Food                []Point
	Snake               []Point
	Direction           Point
	GameOver            bool
	Score               int
	FoodNearby          bool
	FoodCount           int
	EnableShrinkingFood bool
	EnablePortalFood    bool
	ShrinkingFruits     []Point
	ShrinkingFoodCount  int
	DirectionQueue      []Point
	PortalFruits        [][2]Point
	Blocks              []Point
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
		PortalFruits:       [][2]Point{{{snakeHead.X + 7, snakeHead.Y + 3}, {snakeHead.X + 7, snakeHead.Y - 3}}},
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
func getEmptyCells(g *Game) []Point {
	width, height := g.Width, g.Height
	var emptyCells []Point
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			currentPoint := Point{X: x, Y: y}
			if !checkSnake(g, currentPoint) && !checkFood(g, currentPoint) && !checkShrinkingFruit(g, currentPoint) && !checkPortalFoods(g, currentPoint) && !checkBlocks(g, currentPoint) {
				emptyCells = append(emptyCells, currentPoint)
			}
		}
	}
	return emptyCells
}
func checkSnake(g *Game, point Point) bool {
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
	if checkBlocks(g, newHead) {
		g.GameOver = true
	}
	if check {
		g.Score++
		g.placeFood()
		if g.Score%2 == 0 {
			g.placeBlocks()
		}
		g.Food = append(g.Food[:eaten], g.Food[eaten+1:]...)
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
	p, peaten, pcheck := eatenPortalFood(g.PortalFruits, newHead)
	if pcheck {
		g.Score++
		g.placePortalFruits()
		if g.Score%2 == 0 {
			g.placeBlocks()
		}
		newHead := Point{X: g.PortalFruits[peaten][p].X, Y: g.PortalFruits[peaten][p].Y}
		g.Snake = append([]Point{newHead}, g.Snake[1:]...)
		g.PortalFruits = append(g.PortalFruits[:peaten], g.PortalFruits[peaten+1:]...)
	}
	if !check && !pcheck {
		g.Snake = g.Snake[:len(g.Snake)-1]
	}

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
