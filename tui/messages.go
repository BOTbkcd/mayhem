package tui

import tea "github.com/charmbracelet/bubbletea"

type goToMainMsg struct {
	value string
}

func goToMainCmd() tea.Msg {
	return goToMainMsg{}
}

func goToMainWithVal(value string) tea.Cmd {
	return func() tea.Msg {
		return goToMainMsg{value: value}
	}
}

type goToFormMsg struct {
	value interface{}
}

func goToFormWithVal(value interface{}) tea.Cmd {
	return func() tea.Msg {
		return goToFormMsg{value: value}
	}
}

type goToStepsMsg struct {
	value interface{}
}

func goToStepsWithVal(value interface{}) tea.Cmd {
	return func() tea.Msg {
		return goToStepsMsg{value: value}
	}
}
