package tui

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	weekDays = map[int]string{
		0: "Mo",
		1: "Tu",
		2: "We",
		3: "Th",
		4: "Fr",
		5: "Sa",
		6: "Su",
	}
)

type calendar struct {
	selectedDate time.Time
	totalDays    int
	startOffset  int
}

func initializeCalender(selectedDate time.Time) calendar {
	c := calendar{}
	c.selectedDate = selectedDate

	//As per time package Sunday has 0 index, but in our arrangement Sunday appears at the end of the row with 7 index
	offset := int(c.selectedDate.AddDate(0, 0, -c.selectedDate.Day()+1).Weekday()) - int(time.Monday)

	if offset == -1 {
		c.startOffset = 6
	} else {
		c.startOffset = offset
	}

	c.totalDays = c.selectedDate.AddDate(0, 1, -c.selectedDate.Day()).Day()

	return c
}

func (c calendar) Init() tea.Cmd {
	return nil
}

func (c calendar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Right):
			newDate := c.selectedDate.AddDate(0, 0, 1)
			return initializeCalender(newDate), nil

		case key.Matches(msg, Keys.Left):
			newDate := c.selectedDate.AddDate(0, 0, -1)
			return initializeCalender(newDate), nil

		case key.Matches(msg, Keys.Up):
			newDate := c.selectedDate.AddDate(0, 0, -7)
			return initializeCalender(newDate), nil

		case key.Matches(msg, Keys.Down):
			newDate := c.selectedDate.AddDate(0, 0, 7)
			return initializeCalender(newDate), nil

		}
	}
	return c, nil
}

func (c calendar) View() string {
	//Add month + year row
	monthRow := lipgloss.NewStyle().Padding(1, 0).Bold(true).Render(c.renderCalendarMonth())

	//Add weekday row
	weekdayRow := lipgloss.NewStyle().Padding(0, 1).Render(c.renderWeekDays())

	return lipgloss.JoinVertical(lipgloss.Center, monthRow, weekdayRow, c.renderCalender())
}

func (c calendar) renderCalender() string {
	var output []string

	//Calculate ceiling
	rowCount := (c.totalDays + c.startOffset + (7 - 1)) / 7

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		renderedRow := c.renderCalendarRow(rowIndex)
		output = append(output, lipgloss.NewStyle().Padding(0, 1).Render(renderedRow))
	}

	return lipgloss.JoinVertical(lipgloss.Top, output...)
}

func (c calendar) renderCalendarRow(rowIndex int) string {
	rowString := make([]string, 7)

	for colIndex := 1; colIndex <= 7; colIndex++ {
		value := c.getBoxValue(rowIndex, colIndex)

		if c.isCurrentDay(value) {
			rowString = append(rowString, c.renderBox(value, timeFocusColor, true))
		} else {
			color := whiteColor
			if colIndex == 7 {
				color = unfocusedColor
			}

			rowString = append(rowString, c.renderBox(value, color, false))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, rowString...)
}

func (c calendar) renderWeekDays() string {
	rowString := make([]string, 7)

	for i := 0; i < 7; i++ {
		value := weekDays[i]
		color := whiteColor
		if i == 6 {
			color = unfocusedColor
		}

		rowString = append(rowString, c.renderBox(value, color, true))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, rowString...)
}

func (c calendar) renderBox(val string, color lipgloss.Color, border bool) string {
	style := lipgloss.NewStyle().
		Foreground(color).
		BorderForeground(color).
		Padding(0, 2)

	if border {
		style = style.Border(lipgloss.RoundedBorder())
	} else {
		style = style.Border(lipgloss.HiddenBorder())
	}

	return style.Render(val)
}

func (c calendar) getBoxValue(rowIndex int, colIndex int) string {
	value := colIndex + 7*rowIndex - c.startOffset
	if value <= 0 || value > c.totalDays {
		return "  "
	} else {
		return fmt.Sprintf("%02d", value)
	}
}

func (c calendar) renderCalendarMonth() string {
	month := c.selectedDate.Month().String()
	year := strconv.Itoa(c.selectedDate.Year())

	return month + " - " + year
}

func (c calendar) isCurrentDay(boxValue string) bool {
	selectedDay := c.selectedDate.Day()
	boxDay, _ := strconv.Atoi(boxValue)

	if selectedDay == boxDay {
		return true
	} else {
		return false
	}
}
