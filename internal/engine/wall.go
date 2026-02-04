package engine

func (g *Game) checkWallCollision(p Point) bool {
	if p.X < 0 || p.Y < 0 || p.X >= g.Width || p.Y >= g.Height {
		return true
	}
	return false
}
func (g *Game) handleWallCollision(newHead *Point) {
	if newHead.X == -1 {
		newHead.X = g.Width - 1
	}
	if newHead.Y == -1 {
		newHead.Y = g.Height - 1
	}
	if newHead.X == g.Width {
		newHead.X = 0
	}
	if newHead.Y == g.Height {
		newHead.Y = 0
	}
}
