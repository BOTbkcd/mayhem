package tui

import (
	"github.com/BOTbkcd/mayhem/entities"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// textinput.Model doesn't implement tea.Model interface
type stepsEditor struct {
	taskID     uint
	steps      []entities.Step
	textInput  tea.Model
	focusIndex int
	isEditMode bool
}

var stepsEditorKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("'↑/k'", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("'↓/j'", "down"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("'n'", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("'e'", "edit"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("'enter'", "save"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("'tab'", "toggle status"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("'x'", "delete"),
	),
	Return: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("'esc'", "return"),
	),
}

var (
	selectedStepColor = lipgloss.Color("#FFFF00")

	unselectedStepColor = lipgloss.Color("#999999")
)

var (
	newStepPlaceholder = "Enter Step Description"
)

func initializeStepsEditor(steps []entities.Step, taskID uint) tea.Model {
	t := stepsEditor{
		taskID:     taskID,
		steps:      steps,
		focusIndex: 0,
	}

	return t
}

func (m stepsEditor) Init() tea.Cmd {
	return nil
}

func (m stepsEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.isEditMode {
		switch msg := msg.(type) {
		case goToStepsMsg:
			currStep := m.steps[m.focusIndex]
			currStep.Title = msg.value.(string)
			m.steps[m.focusIndex] = currStep.Save()
			m.isEditMode = false
			return m, nil
		}

		// Placing it outside KeyMsg case is required, otherwise messages like textinput's Blink will be lost
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, Keys.Up):
			if m.focusIndex > 0 {
				m.focusIndex--
			}

		case key.Matches(msg, Keys.Down):
			if m.focusIndex < len(m.steps)-1 {
				m.focusIndex++
			}

		case key.Matches(msg, Keys.New):
			newStep := entities.Step{
				TaskID:     m.taskID,
				Title:      "",
				IsFinished: false,
			}

			if len(m.steps) == 0 {
				m.steps = []entities.Step{newStep}
			} else {
				m.steps = append(m.steps[:m.focusIndex+1], m.steps[m.focusIndex:]...)
				m.steps[m.focusIndex+1] = newStep
				m.focusIndex++
			}
			m.isEditMode = true
			m.textInput = initializeTextInput("", newStepPlaceholder, 60, 0, goToStepsWithVal)

		case key.Matches(msg, Keys.Delete):
			if len(m.steps) > 0 {
				step := m.steps[m.focusIndex]
				step.Delete()
				m.steps = append(m.steps[:m.focusIndex], m.steps[m.focusIndex+1:]...)

				if m.focusIndex > 0 {
					m.focusIndex--
				}
			}

		case key.Matches(msg, Keys.Toggle):
			if len(m.steps) > 0 {
				currStep := m.steps[m.focusIndex]
				currStep.IsFinished = !currStep.IsFinished
				currStep.Save()
				m.steps[m.focusIndex] = currStep
			}

		case key.Matches(msg, Keys.Edit):
			if len(m.steps) > 0 {
				m.isEditMode = true
				m.textInput = initializeTextInput(m.steps[m.focusIndex].Title, newStepPlaceholder, 60, 0, goToStepsWithVal)
			}

		case key.Matches(msg, Keys.Enter):
			return m, goToFormWithVal("")
		}
	}
	return m, nil
}

func (m stepsEditor) View() string {
	var res []string

	for i, val := range m.steps {
		var value string

		if i == m.focusIndex && m.isEditMode {
			res = append(res, m.textInput.View())
			continue
		} else if i == m.focusIndex {
			value = "» "
		} else {
			value = "  "
		}

		if val.IsFinished {
			value += lipgloss.NewStyle().Strikethrough(true).Render(val.Title)
			res = append(res, value)
		} else {
			value += val.Title
			res = append(res, value)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, res...)
}
