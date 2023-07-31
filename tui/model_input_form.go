package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/BOTbkcd/mayhem/entities"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type inputForm struct {
	focusIndex    int
	data          entities.Entity
	dataType      string
	fieldMap      map[int]field
	isInvalid     bool
	invalidPrompt string
	isNewTask     bool
	helpKeys      keyMap
}

type field struct {
	name             string
	prompt           string
	model            tea.Model
	isRequired       bool
	nilValue         string
	validationPrompt string
	helpKeys         keyMap
}

var (
	stackFields map[int]field = map[int]field{
		0: {
			name:             "Title",
			prompt:           "Stack Title",
			isRequired:       true,
			nilValue:         "",
			helpKeys:         textInputKeys,
			validationPrompt: "Stack title field can not be empty❗",
		},
	}

	taskFields map[int]field = map[int]field{
		0: {
			name:             "Title",
			prompt:           "Task Title",
			isRequired:       true,
			nilValue:         "",
			helpKeys:         textInputKeys,
			validationPrompt: "Task title field can not be empty❗",
		},
		1: {
			name:     "Description",
			prompt:   "Task Description",
			helpKeys: textAreaKeys,
		},
		2: {
			name:     "Steps",
			prompt:   "Task Steps",
			helpKeys: stepsEditorKeys,
		},
		3: {
			name:     "Priority",
			prompt:   "Task Priority",
			helpKeys: listSelectorKeys,
		},
		4: {
			name:     "Deadline",
			prompt:   "Task Deadline",
			helpKeys: timePickerKeys,
		},
		5: {
			name:     "StartAt",
			prompt:   "Task Start Time",
			helpKeys: timePickerKeys,
		},
		6: {
			name:     "RecurrenceInterval",
			prompt:   "Task Recurrence Interval",
			helpKeys: timePickerKeys,
		},
	}
)

var (
	StackFieldIndex map[string]int = map[string]int{
		"Title": 0,
	}

	TaskFieldIndex map[string]int = map[string]int{
		"Title":              0,
		"Description":        1,
		"Steps":              2,
		"Priority":           3,
		"Deadline":           4,
		"StartAt":            5,
		"RecurrenceInterval": 6,
	}
)

func initializeInput(selectedTable string, data entities.Entity, fieldIndex int) inputForm {
	var m inputForm
	if selectedTable == "stack" {
		m = inputForm{
			data:       data,
			focusIndex: fieldIndex,
			dataType:   "stack",
			fieldMap:   stackFields,
		}

		targetField := m.fieldMap[fieldIndex]
		stack := data.(entities.Stack)

		switch fieldIndex {
		case 0:
			targetField.model = initializeTextInput(stack.Title, "", 20, goToFormWithVal)
		}

		m.helpKeys = targetField.helpKeys
		m.fieldMap[fieldIndex] = targetField

	} else {
		m = inputForm{
			data:       data,
			focusIndex: fieldIndex,
			fieldMap:   taskFields,
			dataType:   "task",
		}

		targetField := m.fieldMap[fieldIndex]
		task := data.(entities.Task)

		switch fieldIndex {
		case 0:
			targetField.model = initializeTextInput(task.Title, "", 60, goToFormWithVal)
		case 1:
			targetField.model = initializeTextArea(task.Description)
		case 2:
			targetField.model = initializeStepsEditor(task.Steps, task.ID)
		case 3:
			opts := []keyVal{
				{val: "0"},
				{val: "1"},
				{val: "2"},
			}
			targetField.model = initializeListSelector(opts, strconv.Itoa(task.Priority), goToFormWithVal)
		case 4:
			if task.Deadline.IsZero() {
				currDate := time.Now().String()[0:10]
				startOfToday, _ := time.Parse(time.DateOnly, currDate)
				targetField.model = initializeTimePicker(startOfToday)
			} else {
				targetField.model = initializeTimePicker(task.Deadline)
			}
		case 5:
			targetField.model = initializeMomentPicker(task.StartTime)
		case 6:
			targetField.model = initializeDurationPicker(task.RecurrenceInterval)
		}
		m.helpKeys = targetField.helpKeys
		m.fieldMap[fieldIndex] = targetField
	}

	return m
}

func (m inputForm) Init() tea.Cmd {
	return nil
}

func (m inputForm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	//Transfer control to selectModel's Update method
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Return):
			if m.focusIndex == TaskFieldIndex["Steps"] {
				//In case of steps editor the steps are saved at the time of editing itself,
				//so returning from steps editor should update the data
				return m, goToMainWithVal("refresh")
			} else {
				return m, goToMainCmd
			}

		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		}

	case goToFormMsg:
		selectedValue := msg.value

		if (m.fieldMap[m.focusIndex].isRequired) && (selectedValue == m.fieldMap[m.focusIndex].nilValue) {
			m.isInvalid = true
			m.invalidPrompt = m.fieldMap[m.focusIndex].validationPrompt
			return m, nil
		} else {
			m.isInvalid = false
		}

		if m.dataType == "stack" {
			stack := m.data.(entities.Stack)

			switch m.focusIndex {
			case 0:
				stack.Title = selectedValue.(string)
			}

			stack.Save()

		} else if m.dataType == "task" {
			task := m.data.(entities.Task)

			switch m.focusIndex {
			case 0:
				task.Title = selectedValue.(string)

				if task.CreatedAt.IsZero() {
					m.isNewTask = true
				}

			case 1:
				task.Description = selectedValue.(string)
			case 2:
				// We save tasks independently (in steps-editor itself) & not as task associations
				return m, goToMainWithVal("refresh")
			case 3:
				task.Priority, _ = strconv.Atoi(selectedValue.(keyVal).val)
			case 4:
				oldDeadline := task.Deadline
				task.Deadline = selectedValue.(time.Time)
				if task.IsRecurring {
					spawnRecurTasks(task, oldDeadline)
				}

			case 5:
				prevTime := task.StartTime
				newTime := selectedValue.(time.Time)
				task.StartTime = time.Date(prevTime.Year(), prevTime.Month(), prevTime.Day(), newTime.Hour(), newTime.Minute(), 0, 0, prevTime.Location())
				spawnRecurTasks(task, task.Deadline)

			case 6:
				task.RecurrenceInterval = selectedValue.(int)
				spawnRecurTasks(task, task.Deadline)
			}

			task = task.Save().(entities.Task)

			if m.isNewTask {
				if task.IsRecurring {
					recurTask := entities.RecurTask{
						StackID:  task.StackID,
						TaskID:   task.ID,
						Deadline: task.StartTime,
					}
					recurTask.Save()
				} else {
					entities.IncPendingCount(task.StackID)
				}
			}
		}

		return m, goToMainWithVal("refresh")
	}

	// Placing it outside KeyMsg case is required, otherwise messages like textinput's Blink will be lost
	var cmd tea.Cmd
	inputField := m.fieldMap[m.focusIndex]
	inputField.model, cmd = m.fieldMap[m.focusIndex].model.Update(msg)
	m.fieldMap[m.focusIndex] = inputField

	return m, cmd
}

func (m inputForm) View() string {
	var b strings.Builder

	//ADD changes for invalid input case

	b.WriteString(highlightedTextStyle.Render(m.fieldMap[m.focusIndex].prompt))

	if m.isInvalid {
		b.WriteString(lipgloss.NewStyle().Foreground(highlightedBackgroundColor).Render("    **" + m.invalidPrompt))
	}

	b.WriteRune('\n')
	b.WriteRune('\n')

	b.WriteString(m.fieldMap[m.focusIndex].model.View())
	b.WriteRune('\n')

	// blurredButton := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	// focusedButton := focusedStyle.Copy().Render("[ Submit ]")

	// var button *string
	// if m.focusIndex == len(m.fieldMap) {
	// 	button = &focusedButton
	// } else {
	// 	button = &blurredButton
	// }

	// fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func spawnRecurTasks(task entities.Task, oldDeadline time.Time) {
	if task.Deadline.Before(time.Now()) {
		return
	}

	r, _ := task.LatestRecurTask()

	//Delete all recur tasks from now
	task.RemoveFutureRecurTasks()

	var startTime time.Time
	t := time.Now()

	if t.Before(oldDeadline) {
		startTime = r.Deadline
	} else {
		startTime = time.Date(t.Year(), t.Month(), t.Day(), task.StartTime.Hour(), task.StartTime.Minute(), 0, 0, task.StartTime.Location())
	}

	for startTime.Compare(task.Deadline) <= 0 {
		recurTask := entities.RecurTask{
			TaskID:     task.ID,
			StackID:    task.StackID,
			IsFinished: false,
			Deadline:   startTime,
		}
		recurTask.Save()

		startTime = startTime.AddDate(0, 0, task.RecurrenceInterval)
	}
}
