package cmd

import (
	"fmt"
	"gosnake/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	W         int
	H         int
	FoodCount int
)

var RootCmd = &cobra.Command{
	Use:   "snake",
	Short: "snake game in your terminal",
	Long:  `snake game in your terminal`,
	Run:   NewGame,
}

func NewGame(_ *cobra.Command, _ []string) {
	p := tea.NewProgram(ui.InitialModel(W, H, FoodCount))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
func init() {
	RootCmd.Flags().IntVarP(&W, "width", "W", 20, "the width of the game")
	RootCmd.Flags().IntVarP(&H, "height", "H", 20, "the height of the game")
	RootCmd.Flags().IntVarP(&FoodCount, "food", "f", 1, "the food count")
}
