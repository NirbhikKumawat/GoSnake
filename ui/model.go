package ui

import (
	"fmt"
	"gosnake/internal/engine"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time
type Model struct {
	game *engine.Game
}

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func InitialModel() Model {
	return Model{
		game: engine.NewGame(20, 20),
	}
}
func (m Model) Init() tea.Cmd {
	return doTick()
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !m.game.GameOver {
			m.game.Move()
			return m, doTick()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.game.Direction = m.game.NextDirection(0, -1)
			return m, nil
		case "down":
			m.game.Direction = m.game.NextDirection(0, 1)
			return m, nil
		case "right":
			m.game.Direction = m.game.NextDirection(1, 0)
			return m, nil
		case "left":
			m.game.Direction = m.game.NextDirection(-1, 0)
			return m, nil
		}
	}
	return m, nil
}
func (m Model) View() string {
	g := m.game
	grid := make([][]string, m.game.Height)
	for i := range grid {
		grid[i] = make([]string, m.game.Width)
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}
	grid[g.Food.Y][g.Food.X] = "@"
	for index, part := range g.Snake {
		char := "o"
		if index == 0 {
			char = "O"
		}
		grid[part.Y][part.X] = char
	}
	var b strings.Builder
	for y := 0; y < m.game.Height; y++ {
		for x := 0; x < m.game.Width; x++ {
			b.WriteString(grid[y][x] + " ")
		}
		b.WriteString("\n")
	}
	if g.GameOver {
		b.WriteString("\nYou are over!\n\n")
	}
	b.WriteString(fmt.Sprintf("\nScore: %d", g.Score))
	return b.String()
}
