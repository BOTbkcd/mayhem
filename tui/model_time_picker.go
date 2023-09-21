package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type timePicker struct {
	currTime         time.Time
	focusIndex       int
	isDurationPicker bool //Show all fields, else only show days
	isMomentPicker   bool //Show only min+hr fields
	dayCount         int  //Used in duration picker mode
}

type timeUnit struct {
	title     string
	tag       string
	charWidth int
}

var timePickerKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("'↑/k'", "increase"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("'↓/j'", "decrease"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("'←/h'", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("'→/l'", "move right"),
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

var timeUnitMap = map[int]timeUnit{
	0: {
		title:     "Hour",
		tag:       "hh",
		charWidth: 2,
	},
	1: {
		title:     "Minute",
		tag:       "mm",
		charWidth: 2,
	},
	2: {
		title:     "Day",
		tag:       "DD",
		charWidth: 2,
	},
	3: {
		title:     "Month",
		tag:       "MM",
		charWidth: 2,
	},
	4: {
		title:     "Year",
		tag:       "YYYY",
		charWidth: 4,
	},
}

func initializeTimePicker(currTime time.Time) tea.Model {
	t := timePicker{
		currTime: currTime,
	}

	return t
}

func initializeDurationPicker(recurrenceInterval int) tea.Model {
	t := timePicker{
		dayCount:         recurrenceInterval,
		isDurationPicker: true,
		focusIndex:       2,
	}

	return t
}

func initializeMomentPicker(currTime time.Time) tea.Model {
	t := timePicker{
		currTime:       currTime,
		isMomentPicker: true,
	}

	return t
}

func (m timePicker) Init() tea.Cmd {
	return nil
}

func (m timePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, Keys.Up):
			if m.isDurationPicker {
				m.dayCount++
			} else {
				switch m.focusIndex {
				case 0:
					hourDuration, _ := time.ParseDuration("60m")
					m.currTime = m.currTime.Add(hourDuration)
				case 1:
					minuteDuration, _ := time.ParseDuration("1m")
					m.currTime = m.currTime.Add(minuteDuration)
				case 2:
					m.currTime = m.currTime.AddDate(0, 0, 1)
				case 3:
					m.currTime = m.currTime.AddDate(0, 1, 0)
				case 4:
					m.currTime = m.currTime.AddDate(1, 0, 0)
				}
			}
			return m, nil

		case key.Matches(msg, Keys.Down):
			if m.isDurationPicker {
				if m.dayCount > 1 {
					m.dayCount--
				}
			} else {
				switch m.focusIndex {
				case 0:
					hourDuration, _ := time.ParseDuration("60m")
					m.currTime = m.currTime.Add(-hourDuration)
				case 1:
					minuteDuration, _ := time.ParseDuration("1m")
					m.currTime = m.currTime.Add(-minuteDuration)
				case 2:
					m.currTime = m.currTime.AddDate(0, 0, -1)
				case 3:
					m.currTime = m.currTime.AddDate(0, -1, 0)
				case 4:
					m.currTime = m.currTime.AddDate(-1, 0, 0)
				}
			}
			return m, nil

		case key.Matches(msg, Keys.Right):
			if !m.isDurationPicker {
				if m.focusIndex < len(timeUnitMap)-1 {
					m.focusIndex++
				}
				return m, nil
			}

		case key.Matches(msg, Keys.Left):
			if !m.isDurationPicker {
				if m.focusIndex > 0 {
					m.focusIndex--
				}
				return m, nil
			}
		case key.Matches(msg, Keys.Enter):
			if m.isDurationPicker {
				return m, goToFormWithVal(m.dayCount)
			} else {
				return m, goToFormWithVal(m.currTime)
			}
		}
	}
	return m, nil
}

func (m timePicker) View() string {
	var timeUnitLabel string
	var timeValue string

	if m.isDurationPicker {
		return lipgloss.JoinHorizontal(lipgloss.Center, m.renderUnitCol(2, m.dayCount), " Day(s)")
	} else if m.isMomentPicker {
		timeUnitLabel = lipgloss.JoinHorizontal(lipgloss.Center,
			m.renderUnitTag(0),
			" ",
			m.renderUnitTag(1),
			" ",
			"  ",
		)

		timeValue = lipgloss.JoinHorizontal(lipgloss.Center,
			m.renderUnitCol(0, formatHour(m.currTime.Hour())),
			":",
			m.renderUnitCol(1, m.currTime.Minute()),
			" ",
			renderMidDayInfo(m.currTime.Hour()),
		)

		return lipgloss.JoinVertical(lipgloss.Center,
			timeValue,
			timeUnitLabel,
		)

	} else {
		//Empty spaces are added to align the label and value rows
		timeUnitLabel = lipgloss.JoinHorizontal(lipgloss.Center,
			m.renderUnitTag(0),
			" ",
			m.renderUnitTag(1),
			" ",
			"  ",
			"   ",
			m.renderUnitTag(2),
			" ",
			m.renderUnitTag(3),
			" ",
			m.renderUnitTag(4),
		)

		timeValue = lipgloss.JoinHorizontal(lipgloss.Center,
			m.renderUnitCol(0, formatHour(m.currTime.Hour())),
			":",
			m.renderUnitCol(1, m.currTime.Minute()),
			" ",
			renderMidDayInfo(m.currTime.Hour()),
			"   ",
			m.renderUnitCol(2, m.currTime.Day()),
			"-",
			m.renderUnitCol(3, int(m.currTime.Month())),
			"-",
			m.renderUnitCol(4, m.currTime.Year()))

		return lipgloss.JoinVertical(lipgloss.Center,
			timeValue,
			timeUnitLabel,
		)
	}

}

func (m timePicker) renderUnitCol(index int, val int) string {
	value := fmt.Sprintf("%0*d", timeUnitMap[index].charWidth, val)

	var color lipgloss.Color
	if m.focusIndex == index {
		color = timeFocusColor
	} else {
		color = unfocusedColor
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		BorderForeground(color).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	return style.Render(value)
}

func (m timePicker) renderUnitTag(index int) string {
	value := timeUnitMap[index].tag

	var color lipgloss.Color
	if m.focusIndex == index {
		color = timeFocusColor
	} else {
		color = unfocusedColor
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		Padding(0, 2)

	return style.Render(value)
}

func renderMidDayInfo(hours int) string {
	if isBeforeMidDay(hours) {
		return "am"
	} else {
		return "pm"
	}
}

// Adjust Hour value to 12 hour clock format
func formatHour(value int) int {
	if value > 12 {
		return value - 12
	} else {
		return value
	}
}

func isBeforeMidDay(value int) bool {
	if value >= 12 {
		return false
	} else {
		return true
	}
}
