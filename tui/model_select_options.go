package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type selectModel struct {
	options    []string
	focusIndex int
	maxIndex   int
}

var selectOptionKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("'↑/k'", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("'↓/j'", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("'enter'", "save"),
	),
	Return: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("'esc'", "return"),
	),
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func initialSelectModel(options []string, selectedVal string) tea.Model {
	// Takes care of default case where index should be 0
	var selectedIndex int

	for i, val := range options {
		if val == selectedVal {
			selectedIndex = i
			break
		}
	}

	m := selectModel{focusIndex: selectedIndex, maxIndex: len(options) - 1, options: options}

	return m
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, goToFormWithVal(m.options[m.focusIndex])
		case "up", "k":
			if m.focusIndex > 0 {
				m.focusIndex--
			} else {
				m.focusIndex = m.maxIndex
				return m, nil
			}
		case "down", "j":
			if m.focusIndex < m.maxIndex {
				m.focusIndex++
			} else {
				m.focusIndex = 0
				return m, nil
			}
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	var res []string

	for i, val := range m.options {
		var value string

		if i == m.focusIndex {
			value = lipgloss.NewStyle().Foreground(inputFormColor).Bold(true).Render("» " + val)
		} else {
			value = lipgloss.NewStyle().Foreground(inputFormColor).Bold(true).Render("  " + val)
		}

		res = append(res, value)
	}

	return lipgloss.JoinVertical(lipgloss.Left, res...)
}
