package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listSelector struct {
	options    []keyVal
	focusIndex int
	maxIndex   int
	responder  func(interface{}) tea.Cmd
}

type keyVal struct {
	key uint
	val string
}

var listSelectorKeys = keyMap{
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

func (m listSelector) Init() tea.Cmd {
	return nil
}

func initializeListSelector(options []keyVal, selectedVal string, responder func(interface{}) tea.Cmd) tea.Model {
	// Takes care of default case where index should be 0
	var selectedIndex int

	for i, item := range options {
		if item.val == selectedVal {
			selectedIndex = i
			break
		}
	}

	m := listSelector{
		focusIndex: selectedIndex,
		maxIndex:   len(options) - 1,
		options:    options,
		responder:  responder,
	}

	return m
}

func (m listSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Return):
			return m, goToMainWithVal(keyVal{})

		case key.Matches(msg, Keys.Quit, Keys.Exit):
			return m, tea.Quit

		case key.Matches(msg, Keys.Enter):
			return m, m.responder(m.options[m.focusIndex])

		case key.Matches(msg, Keys.Up):
			if m.focusIndex > 0 {
				m.focusIndex--
			} else {
				m.focusIndex = m.maxIndex
				return m, nil
			}

		case key.Matches(msg, Keys.Down):
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

func (m listSelector) View() string {
	var res []string

	for i, item := range m.options {
		var value string

		if i == m.focusIndex {
			value = lipgloss.NewStyle().Foreground(inputFormColor).Bold(true).Render("» " + item.val)
		} else {
			value = lipgloss.NewStyle().Foreground(inputFormColor).Bold(true).Render("  " + item.val)
		}

		res = append(res, value)
	}

	return lipgloss.JoinVertical(lipgloss.Left, res...)
}
