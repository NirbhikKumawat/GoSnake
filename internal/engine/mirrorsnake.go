package engine

func checkMirrorSnake(g *Game, point Point) bool {
	for _, snakeHead := range g.MirrorSnake {
		if snakeHead.X == point.X && snakeHead.Y == point.Y {
			return true
		}
	}
	return false
}
func (g *Game) mirrorCollision(p1, p2 Point) bool {
	if p1.X == p2.X && p1.Y == p2.Y {
		return true
	}
	return false
}
