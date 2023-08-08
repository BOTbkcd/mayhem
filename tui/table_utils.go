package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BOTbkcd/mayhem/entities"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
)

var stackKeys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("'n'", "new stack 🌟"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("'e'", "edit 📝"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("'x'", "delete 🗑"),
	),
}

var taskKeys = keyMap{
	Toggle: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("'tab'", "check/uncheck 🔄"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("'n'", "new task 🌟"),
	),
	NewRecur: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("'r'", "new recurring task 🌟"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("'e'", "edit 📝"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("'x'", "delete 🗑"),
	),
	Move: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("'m'", "change stack 📤"),
	),
}

var tableNavigationKeys = keyMap{
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
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("'←/h'", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("'→/l'", "right"),
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

var taskFinishStatus = map[uint]bool{}
var recurDeadlines = map[uint]time.Time{}

func stackColumns() []table.Column {
	return []table.Column{
		{Title: "       Stacks 🗃", Width: 20},
		{Title: "", Width: 2},
	}
}

func taskColumns() []table.Column {
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "           Tasks 📄", Width: 30},
		{Title: "     Deadline 🕑", Width: 20},
		{Title: "Priority", Width: 8},
	}
}

func stackRows(stacks []entities.Stack) []table.Row {
	rows := make([]table.Row, len(stacks))

	sortStacks(stacks)

	for i, val := range stacks {
		row := []string{
			val.Title,
			incompleteTaskTag(val.PendingTaskCount + val.PendingRecurringCount()),
		}
		rows[i] = row
	}
	return rows
}

func taskRows(tasks []entities.Task) []table.Row {
	rows := make([]table.Row, len(tasks))

	// We perform this step earlier since we need the deadline & finish status data before sorting
	for _, val := range tasks {
		if val.IsRecurring {
			r, count := val.LatestRecurTask()
			if count > 0 {
				recurDeadlines[val.ID] = r.Deadline
				taskFinishStatus[val.ID] = r.IsFinished
			}
		} else {
			taskFinishStatus[val.ID] = val.IsFinished
		}
	}

	sortTasks(tasks)

	var prefix string
	var deadline string

	for i, val := range tasks {
		if val.IsRecurring {
			deadline = formatTime(recurDeadlines[val.ID], true)
			// prefix = "𝑹"
			// prefix = "🅁 Ⓡ 🄬"
		} else {
			deadline = formatTime(val.Deadline, true)
		}

		if taskFinishStatus[val.ID] {
			prefix = "✘"
		} else {
			prefix = "▢"
		}

		row := []string{
			prefix,
			val.Title,
			deadline,
			"   " + strconv.Itoa(val.Priority),
		}

		if val.IsRecurring {
			row[3] = row[3] + "  📌"
		}

		rows[i] = row
	}

	return rows
}

func sortStacks(s []entities.Stack) {
	//Alphabetically sort by stack title
	sort.Slice(s, func(i, j int) bool {
		return strings.ToLower(s[i].Title) < strings.ToLower(s[j].Title)
	})
}

func sortTasks(t []entities.Task) {
	//Sort by finish status, then deadline, then priority, then title
	sort.Slice(t, func(i, j int) bool {
		if taskFinishStatus[t[i].ID] == taskFinishStatus[t[j].ID] {
			var deadline_i time.Time
			if t[i].IsRecurring {
				deadline_i = recurDeadlines[t[i].ID]
			} else {
				deadline_i = t[i].Deadline
			}

			var deadline_j time.Time
			if t[j].IsRecurring {
				deadline_j = recurDeadlines[t[j].ID]
			} else {
				deadline_j = t[j].Deadline
			}

			if deadline_i.Equal(deadline_j) {
				if t[i].Priority == t[j].Priority {
					return strings.ToLower(t[i].Title) < strings.ToLower(t[j].Title)
				}
				return t[i].Priority > t[j].Priority
			}

			if deadline_i.IsZero() {
				return false
			}

			if deadline_j.IsZero() {
				return true
			}

			return deadline_i.Before(deadline_j)
		} else {
			return !taskFinishStatus[t[i].ID]
		}
	})
}

func buildTable(columns []table.Column, tableType string) table.Model {
	t := table.New(
		table.WithHeight(tableViewHeight),
		table.WithColumns(columns),
		table.WithKeyMap(table.DefaultKeyMap()),
		// table.WithFocused(true),
	)

	s := getTableStyle(tableType)
	t.SetStyles(s)

	return t
}

func formatTime(time time.Time, fullDate bool) string {
	if time.IsZero() {
		return fmt.Sprintf("%10s", dash)
	}

	year := fmt.Sprintf("%04d", time.Year())
	month := fmt.Sprintf("%02d", int(time.Month()))
	days := fmt.Sprintf("%02d", time.Day())
	hours := fmt.Sprintf("%02d", formatHour(time.Hour()))
	minutes := fmt.Sprintf("%02d", time.Minute())
	midDayInfo := renderMidDayInfo(time.Hour())

	if fullDate {
		return days + "-" + month + "-" + year + "  " + hours + ":" + minutes + " " + midDayInfo
	} else {
		return hours + ":" + minutes + " " + midDayInfo
	}

}

func getEmptyTaskView() string {
	return getEmptyTaskStyle().Render("Press either '→' or 'l' key to explore this stack")
}

func getEmptyDetailsView() string {
	return getEmptyDetailsStyle().Render("Press either '→' or 'l' key to see task details")
}

func incompleteTaskTag(count int) string {
	if count > 0 && count <= 10 {
		return " " + string(rune('➊'+count-1))
	} else if count > 10 {
		return "+➓"
	} else {
		return ""
	}
}

func min(x int, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func findStackIndex(arr []entities.Stack, id uint) int {
	for i, val := range arr {
		if val.ID == id {
			return i
		}
	}
	return -1
}

func findTaskIndex(arr []entities.Task, id uint) int {
	for i, val := range arr {
		if val.ID == id {
			return i
		}
	}
	return -1
}
