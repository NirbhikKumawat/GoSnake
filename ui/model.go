package ui

import (
	"fmt"
	"gosnake/internal/engine"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time
type Model struct {
	game *engine.Game
}

var (
	foodStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B0000"))
	snakeHeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	snakeStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	nearbyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#00008B"))
	sfoodStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
)

func doTick() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func InitialModel(w, h, f, sf int) Model {
	if w < 5 || h < 5 {
		w = 20
		h = 20
	}
	return Model{
		game: engine.NewGame(w, h, f, sf),
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
		case "r":
			if m.game.GameOver {
				m.game = engine.NewGame(m.game.Width, m.game.Height, m.game.FoodCount, m.game.ShrinkingFoodCount)
				return m, doTick()
			}
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
	for _, fruit := range g.Food {
		grid[fruit.Y][fruit.X] = foodStyle.Render("⬛")
	}
	for _, sfruit := range g.ShrinkingFruits {
		grid[sfruit.Y][sfruit.X] = sfoodStyle.Render("⬛")
	}

	for index, part := range g.Snake {
		char := snakeStyle.Render("⬛")
		if index == 0 {
			if !g.FoodNearby {
				char = snakeHeaderStyle.Render("⬛")
			} else {
				char = nearbyStyle.Render("⬛")
			}
		}
		grid[part.Y][part.X] = char
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("\nScore: %d \n", g.Score))
	for y := 0; y < m.game.Height; y++ {
		for x := 0; x < m.game.Width; x++ {
			b.WriteString(grid[y][x] + " ")
		}
		b.WriteString("\n")
	}
	if g.GameOver {
		b.WriteString("\nGame Over!\n")
		b.WriteString("Press r to restart the game.\n")
	}
	return b.String()
}
