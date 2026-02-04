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
	DirectionQueue      []Point
	GameOver            bool
	Score               int
	FoodNearby          bool
	FoodCount           int
	WallCollision       bool
	SnakeCollision      bool
	ReverseDirection    bool
	EnableMirrorSnake   bool
	MirrorSnake         []Point
	EnableShrinkingFood bool
	ShrinkingFruits     []Point
	ShrinkingFoodCount  int
	EnablePortalFood    bool
	PortalFruits        [][2]Point
	SpawnBlocks         bool
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
func NewGame(width, height, count, scount int, wallCollision, snakeCollision, reverseDirection, spawnBlocks, mirrorSnake bool) *Game {
	snakeHead := Point{X: width/2 - 5, Y: height / 2}
	mirrorSnakeHead := Point{X: width/2 + 5, Y: height / 2}
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
		FoodCount:         count,
		Food:              []Point{{snakeHead.X + 7, snakeHead.Y}},
		FoodNearby:        true,
		WallCollision:     wallCollision,
		SnakeCollision:    snakeCollision,
		ReverseDirection:  reverseDirection,
		SpawnBlocks:       spawnBlocks,
		EnableMirrorSnake: mirrorSnake,
	}
	if scount > 0 {
		g.EnableShrinkingFood = true
	}
	if g.EnablePortalFood {
		g.PortalFruits = [][2]Point{{{snakeHead.X + 7, snakeHead.Y + 3}, {snakeHead.X + 7, snakeHead.Y - 3}}}
	}
	if g.EnableShrinkingFood {
		g.ShrinkingFoodCount = scount
		g.ShrinkingFruits = []Point{{snakeHead.X + 7, snakeHead.Y + 5}}
		g.initialPlaceShrinkingFoods()
	}
	if g.EnableMirrorSnake {
		g.MirrorSnake = []Point{
			mirrorSnakeHead,
			{
				X: mirrorSnakeHead.X + 1,
				Y: mirrorSnakeHead.Y,
			},
			{
				X: mirrorSnakeHead.X + 2,
				Y: mirrorSnakeHead.Y,
			},
		}
	}
	g.initialPlaceFoods()
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
	if g.EnableMirrorSnake {
		if checkMirrorSnake(g, point) {
			return true
		}
	}
	return false
}
func (g *Game) Move() {
	if g.GameOver {
		return
	}
	if len(g.Snake) == 0 {
		g.GameOver = true
		return
	}
	newHead := Point{X: g.Snake[0].X + g.Direction.X, Y: g.Snake[0].Y + g.Direction.Y}
	var mirrorNewHead Point
	if g.EnableMirrorSnake {
		mirrorNewHead = Point{X: g.MirrorSnake[0].X - g.Direction.X, Y: g.MirrorSnake[0].Y - g.Direction.Y}
	}
	if g.selfCollision(newHead) && !g.SnakeCollision {
		g.GameOver = true
		return
	}
	if g.EnableMirrorSnake {
		if g.mirrorCollision(newHead, mirrorNewHead) {
			g.GameOver = true
			return
		}
	}
	if !g.WallCollision && g.checkWallCollision(newHead) {
		g.GameOver = true
		return
	}
	if g.WallCollision {
		g.handleWallCollision(&newHead)
	}
	if g.checkFoodNearby(newHead) {
		g.FoodNearby = true
	} else {
		if g.EnableMirrorSnake {
			if g.checkFoodNearby(mirrorNewHead) {
				g.FoodNearby = true
			}
		}
		g.FoodNearby = false
	}

	g.Snake = append([]Point{newHead}, g.Snake...)
	g.MirrorSnake = append([]Point{mirrorNewHead}, g.MirrorSnake...)

	eaten, check := eatenFood(g.Food, newHead)
	if checkBlocks(g, newHead) {
		g.GameOver = true
		return
	}
	if check {
		g.handleEatenFood(eaten)
	}

	if g.EnableMirrorSnake {
		meaten, mcheck := eatenFood(g.Food, mirrorNewHead)
		if mcheck {
			g.Score++
			g.placeFood()
			g.Food = append(g.Food[:meaten], g.Food[meaten+1:]...)
			check = mcheck
		}
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
		return
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
		if g.EnableMirrorSnake {
			g.MirrorSnake = g.MirrorSnake[:len(g.MirrorSnake)-1]
		}
	}

}
func (g *Game) UpdateDirectionQueue(x, y int) {
	directionx := g.Direction.X
	directiony := g.Direction.Y
	if directionx+x != 0 || directiony+y != 0 && g.Direction != g.Direction && len(g.DirectionQueue) < 3 {
		g.DirectionQueue = append(g.DirectionQueue, Point{X: x, Y: y})
	}
}
func (g *Game) selfCollision(p Point) bool {
	for i := 1; i < len(g.Snake); i++ {
		if p.X == g.Snake[i].X && p.Y == g.Snake[i].Y {
			return true
		}
	}
	return false
}
func (g *Game) getDirection() Point {
	n := len(g.Snake)
	if n >= 2 {
		l := g.Snake[n-1]
		ll := g.Snake[n-2]
		x := l.X - ll.X
		y := l.Y - ll.Y
		return Point{X: x, Y: y}
	} else {
		return Point{X: -g.Direction.X, Y: -g.Direction.Y}
	}
}
