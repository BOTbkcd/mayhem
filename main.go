package main

import (
	"log"

	entities "github.com/BOTbkcd/mayhem/entities"
	tui "github.com/BOTbkcd/mayhem/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	entities.InitializeDB()

	model := tui.InitializeMainModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	// f, err := tea.LogToFile("debug.log", "debug")
	// if err != nil {
	// 	fmt.Println("fatal:", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()

	if _, err := p.Run(); err != nil {
		log.Fatal("Error encountered while running the program:", err)
	}
}
