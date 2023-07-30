package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// textinput.Model doesn't implement tea.Model interface
type deleteConfirmation struct {
}

func initializeDeleteConfirmation() tea.Model {
	m := deleteConfirmation{}

	return m
}

func (m deleteConfirmation) Init() tea.Cmd {
	return textinput.Blink
}

func (m deleteConfirmation) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, Keys.Return):
			return m, goToMainCmd

		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit

		default:
			if msg.String() == "y" || msg.String() == "Y" {
				return m, goToMainWithVal("y")
			} else {
				return m, goToMainWithVal("")
			}
		}
	}
	return m, nil
}

func (m deleteConfirmation) View() string {
	// Can't just render textinput.Value(), otherwise cursor blinking wouldn't work
	return lipgloss.NewStyle().Foreground(highlightedBackgroundColor).PaddingTop(1).Render("Do you wish to proceed with deletion? (y/n): ")
}
