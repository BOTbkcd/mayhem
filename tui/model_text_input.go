package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// textinput.Model doesn't implement tea.Model interface
type textInput struct {
	input textinput.Model
	//Since textinput field can be used in multiple places,
	//responder is required to determine the receiver of the message emitted by textinput field
	responder func(interface{}) tea.Cmd
}

var textInputKeys = keyMap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("'enter'", "save"),
	),
	Return: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("'esc'", "return"),
	),
}

func initializeTextInput(value string, placeholder string, charLimit int, width int, responder func(interface{}) tea.Cmd) tea.Model {
	t := textinput.New()
	t.SetValue(value)

	t.Cursor.Style = textInputStyle
	t.CharLimit = charLimit
	t.Focus()
	t.PromptStyle = textInputStyle
	t.TextStyle = textInputStyle
	t.Placeholder = placeholder
	t.Width = width
	t.EchoCharacter = 'x'

	m := textInput{
		input:     t,
		responder: responder,
	}

	return m
}

func (m textInput) Init() tea.Cmd {
	return textinput.Blink
}

func (m textInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Enter):
			return m, m.responder(m.input.Value())
		}
	}

	// Placing it outside KeyMsg case is required, otherwise messages like textinput's Blink will be lost
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m textInput) View() string {
	// Can't just render textinput.Value(), otherwise cursor blinking wouldn't work
	return m.input.View()
}
