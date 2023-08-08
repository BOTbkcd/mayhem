package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BOTbkcd/mayhem/entities"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type detailsBox struct {
	taskData             entities.Task
	viewport             viewport.Model
	preserveOffset       bool
	oldViewportOffset    int
	focusIndex           int
	isBoxFocused         bool
	isRecurrenceDuration bool
	scrollData           scrollData
}

type scrollData struct {
	title              int
	description        int
	steps              int
	priority           int
	deadline           int
	startTime          int
	recurrenceInterval int
}

var taskDetailsKeys = keyMap{
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("'e'", "edit field 📝"),
	),
	// Toggle: key.NewBinding(
	// 	key.WithKeys("tab"),
	// 	key.WithHelp("'tab'", "next 🔽"),
	// ),
	// ReverseToggle: key.NewBinding(
	// 	key.WithKeys("shift+tab"),
	// 	key.WithHelp("'shift+tab'", "previous 🔼"),
	// ),
}

var detailsNavigationKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("'↑/k'", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("'↓/j'", "down"),
	),
	GotoTop: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("'g'", "jump to top"),
	),
	GotoBottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("'G'", "jump to bottom"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("'?'", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("'q'", "quit"),
	),
}

func (m *detailsBox) buildDetailsBox(data entities.Task, preserveOffset bool) {
	m.taskData = data

	//We want to preserve offset when we return to same details view after editing any field
	//But when going from one task to another, we want to reset the view
	m.preserveOffset = preserveOffset
	m.oldViewportOffset = m.viewport.YOffset
	m.viewport = viewport.New(getDetailsBoxWidth(), tableViewHeight)
	m.renderContent()
}

func (m detailsBox) Init() tea.Cmd {
	return nil
}

func (m detailsBox) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.isBoxFocused {
		return m, nil
	}

	m.viewport.Width = getDetailsBoxWidth()

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, Keys.Up):
			var scrollDistance int
			switch m.focusIndex {
			case 0:
				m.viewport.GotoBottom()
				m.End()
				return m, nil
			case 1:
				scrollDistance = m.scrollData.description
				m.Previous()
			case 2:
				scrollDistance = m.scrollData.steps
				m.Previous()
			case 3:
				scrollDistance = m.scrollData.priority
				m.Previous()
			case 4:
				if m.taskData.IsRecurring {
					scrollDistance = m.scrollData.deadline
				}
				m.Previous()
			case 5:
				scrollDistance = m.scrollData.startTime
				m.Previous()
			case 6:
				// scrollDistance = m.scrollData.recurrenceInterval
				m.Previous()
			}

			m.viewport.LineUp(scrollDistance)

		case key.Matches(msg, Keys.Down):
			var scrollDistance int
			switch m.focusIndex {
			case 0:
				// scrollDistance = m.scrollData.title
				m.Next()
			case 1:
				scrollDistance = m.scrollData.description
				m.Next()
			case 2:
				scrollDistance = m.scrollData.steps
				m.Next()
			case 3:
				scrollDistance = m.scrollData.priority
				m.Next()
			case 4:
				scrollDistance = m.scrollData.deadline
				if m.taskData.IsRecurring {
					m.Next()
				} else {
					m.viewport.GotoTop()
					m.Start()
					return m, nil
				}
			case 5:
				scrollDistance = m.scrollData.startTime
				m.Next()
			case 6:
				m.viewport.GotoTop()
				m.Start()
				return m, nil
			}

			m.viewport.LineDown(scrollDistance)

		case key.Matches(msg, Keys.GotoTop):
			m.viewport.GotoTop()
			m.Start()

		case key.Matches(msg, Keys.GotoBottom):
			m.viewport.GotoBottom()
			m.End()

			// case key.Matches(msg, Keys.Toggle):
			// 	m.Next()

			// case key.Matches(msg, Keys.ReverseToggle):
			// 	m.Previous()
		}
	}
	return m, nil
}

func (m detailsBox) View() string {
	return lipgloss.JoinVertical(lipgloss.Center, getDetailsBoxStyle().Render(m.viewport.View()), m.footerView())
}

func (m *detailsBox) Focus() {
	m.isBoxFocused = true
}

func (m *detailsBox) Blur() {
	m.isBoxFocused = false
}

func (m detailsBox) Focused() bool {
	return m.isBoxFocused
}

func (m *detailsBox) Next() {
	var length int
	if m.taskData.IsRecurring {
		length = 7
	} else {
		length = 5
	}
	m.focusIndex = (m.focusIndex + 1) % length
	m.renderContent()
}

func (m *detailsBox) End() {
	if m.taskData.IsRecurring {
		m.focusIndex = 6
	} else {
		m.focusIndex = 4
	}
	m.renderContent()
}

func (m *detailsBox) Previous() {
	var length int
	if m.taskData.IsRecurring {
		length = 7
	} else {
		length = 5
	}
	val := (m.focusIndex - 1) % length
	if val < 0 {
		val = val + length
	}
	m.focusIndex = val
	m.renderContent()
}

func (m *detailsBox) Start() {
	m.focusIndex = 0
	m.renderContent()
}

func (m *detailsBox) renderContent() {
	var content []string

	if m.taskData.IsRecurring {
		content = []string{
			m.titleBlock(),
			m.descriptionBlock(),
			m.stepsBlock(),
			m.priorityBlock(),
			m.deadlineBlock(),
			m.startTimeBlock(),
			m.recurrenceIntervalBlock(),
		}
	} else {
		content = []string{
			m.titleBlock(),
			m.descriptionBlock(),
			m.stepsBlock(),
			m.priorityBlock(),
			m.deadlineBlock(),
		}
	}

	view := lipgloss.JoinVertical(lipgloss.Left, content...)
	m.viewport.SetContent(view)
	if m.preserveOffset {
		m.viewport.SetYOffset(m.oldViewportOffset)
		m.preserveOffset = false
	}
}
func (m *detailsBox) titleBlock() string {
	var b strings.Builder

	b.WriteString(highlightedTextStyle.Render("Title:"))
	b.WriteString("\n\n")
	b.WriteString(m.taskData.Title)

	isFocused := (m.focusIndex == 0)

	data := getItemContainerStyle(isFocused).Render(getDetailsItemStyle(isFocused).PaddingTop(0).Render(b.String()))
	m.scrollData.title = lipgloss.Height(data)
	return data
}

func (m *detailsBox) descriptionBlock() string {
	var b strings.Builder

	b.WriteString(highlightedTextStyle.Render("Description:"))
	b.WriteString("\n\n")
	if m.taskData.Description == "" {
		b.WriteString(dash)
	} else {
		b.WriteString(m.taskData.Description)
	}

	isFocused := (m.focusIndex == 1)

	data := getItemContainerStyle(isFocused).Render(getDetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.description = lipgloss.Height(data)
	return data
}

func (m *detailsBox) stepsBlock() string {
	var b strings.Builder

	b.WriteString(highlightedTextStyle.Render("Steps:"))
	b.WriteString("\n\n")
	b.WriteString(renderSteps(m.taskData.Steps))

	isFocused := (m.focusIndex == 2)

	data := getItemContainerStyle(isFocused).Render(getDetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.steps = lipgloss.Height(data)
	return data
}

func (m *detailsBox) priorityBlock() string {
	var b strings.Builder

	b.WriteString(highlightedTextStyle.Render("Priority:"))
	b.WriteString("\n\n")
	b.WriteString(strconv.Itoa(m.taskData.Priority))

	isFocused := (m.focusIndex == 3)

	data := getItemContainerStyle(isFocused).Render(getDetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.priority = lipgloss.Height(data)
	return data
}

func (m *detailsBox) deadlineBlock() string {
	var b strings.Builder

	b.WriteString(highlightedTextStyle.Render("Deadline:"))
	b.WriteString("\n\n")
	if m.taskData.Deadline.IsZero() {
		b.WriteString("Not Scheduled")
	} else {
		b.WriteString(formatTime(m.taskData.Deadline, true))
	}

	isFocused := (m.focusIndex == 4)

	data := getItemContainerStyle(isFocused).Render(getDetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.deadline = lipgloss.Height(data)
	return data
}

func (m *detailsBox) startTimeBlock() string {
	var b strings.Builder

	b.WriteString(highlightedTextStyle.Render("Due Time:"))
	b.WriteString("\n\n")
	b.WriteString(formatTime(m.taskData.StartTime, false))

	isFocused := (m.focusIndex == 5)

	data := getItemContainerStyle(isFocused).Render(getDetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.startTime = lipgloss.Height(data)
	return data
}

func (m *detailsBox) recurrenceIntervalBlock() string {
	var b strings.Builder

	b.WriteString(highlightedTextStyle.Render("Recurrence Interval:"))
	b.WriteString("\n\n")
	b.WriteString(strconv.Itoa(m.taskData.RecurrenceInterval) + " day(s)")

	isFocused := (m.focusIndex == 6)

	data := getItemContainerStyle(isFocused).Render(getDetailsItemStyle(isFocused).Render(b.String()))
	m.scrollData.recurrenceInterval = lipgloss.Height(data)
	return data
}

func (m *detailsBox) footerView() string {
	scrollInfoStyle := footerContainerStyle.Copy().
		Width(m.viewport.Width).
		Align(lipgloss.Right)

	info := footerInfoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	return scrollInfoStyle.Render(info)
}

func renderSteps(steps []entities.Step) string {
	var res []string

	if len(steps) == 0 {
		return dash
	}
	for _, val := range steps {
		if val.IsFinished {
			value := lipgloss.JoinHorizontal(
				lipgloss.Center,
				boxedValue("✘"),
				" ",
				val.Title,
			)
			res = append(res, value)
		} else {
			value := lipgloss.JoinHorizontal(
				lipgloss.Center,
				boxedValue(" "),
				" ",
				val.Title,
			)
			res = append(res, value)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, res...)
}

func boxedValue(value string) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render(value)

	return style
}
