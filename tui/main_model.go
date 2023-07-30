package tui

import (
	"fmt"
	"time"

	"github.com/BOTbkcd/mayhem/entities"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	data               []entities.Stack
	stackTable         table.Model
	taskTable          table.Model
	taskDetails        detailsBox
	help               helpModel
	input              inputForm
	showTasks          bool
	showDetails        bool
	showInput          bool
	showHelp           bool
	deleteConfirmation tea.Model
	showDelete         bool
	navigationKeys     keyMap
	preInputFocus      string //useful for reverting back when input box is closed
	firstRender        bool
	prevState          preserveState
}

type preserveState struct {
	retainState bool
	stackID     uint
	taskID      uint
}

func InitializeMainModel() *model {
	stacks, _ := entities.FetchAllStacks()

	m := &model{
		stackTable:     buildTable(stackColumns(), "stack"),
		taskTable:      buildTable(taskColumns(), "task"),
		taskDetails:    detailsBox{}, // we can't build the details box at this stage since we need both stack & task indices for that
		data:           stacks,
		help:           initializeHelp(stackKeys),
		navigationKeys: tableNavigationKeys,
		showHelp:       true,
	}

	m.stackTable.Focus()
	m.taskTable.Blur()
	m.taskDetails.Blur()
	return m
}

func (m *model) Init() tea.Cmd {
	m.firstRender = true
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//Transfer control to inputForm's Update method
	if m.showInput {
		switch msg := msg.(type) {

		case goToMainMsg:
			m.input = inputForm{}
			m.showInput = false

			if msg.value == "refresh" {
				m.preserveState()
				m.refreshData()
			}

			if m.preInputFocus == "stack" {
				m.stackTable.Focus()
				m.help = initializeHelp(stackKeys)
				m.navigationKeys = tableNavigationKeys
			} else if m.preInputFocus == "task" {
				m.taskTable.Focus()
				m.help = initializeHelp(taskKeys)
				m.navigationKeys = tableNavigationKeys
			} else if m.preInputFocus == "detail" {
				m.taskDetails.Focus()
				m.navigationKeys = detailsNavigationKeys
			}
			return m, nil

		default:
			inp, cmd := m.input.Update(msg)
			t, _ := inp.(inputForm)
			m.input = t

			return m, cmd
		}
	}

	//Transfer control to delete confirmation model
	if m.showDelete {
		switch msg := msg.(type) {

		case goToMainMsg:
			m.showDelete = false

			if msg.value == "y" {
				switch m.preInputFocus {
				case "stack":
					stackIndex := m.stackTable.Cursor()
					currStack := m.data[stackIndex]

					if stackIndex == len(m.stackTable.Rows())-1 {
						m.stackTable.SetCursor(stackIndex - 1)
					}

					currStack.Delete()
					m.showTasks = false
					m.showDetails = false
					m.refreshData()
					return m, nil

				case "task":
					stackIndex := m.stackTable.Cursor()
					taskIndex := m.taskTable.Cursor()

					var currTask entities.Task
					if len(m.data[stackIndex].Tasks) > 0 {
						currTask = m.data[stackIndex].Tasks[taskIndex]

						if currTask.IsRecurring {

						} else {
							if !currTask.IsFinished {
								stack := m.data[stackIndex]
								stack.PendingTaskCount--
								stack.Save()
							}
						}
						if taskIndex == len(m.taskTable.Rows())-1 {
							m.taskTable.SetCursor(taskIndex - 1)
						}
						currTask.Delete()
						m.refreshData()
						return m, nil
					}
				}
			}

		default:
			inp, cmd := m.deleteConfirmation.Update(msg)
			t, _ := inp.(deleteConfirmation)
			m.deleteConfirmation = t

			return m, cmd
		}
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		// Inter-table navigation
		case key.Matches(msg, Keys.Left):
			if m.stackTable.Focused() {
				if m.showDetails {
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = initializeHelp(taskDetailsKeys)
					m.navigationKeys = detailsNavigationKeys

				}
			} else if m.taskTable.Focused() {
				m.stackTable.Focus()
				m.taskTable.Blur()
				m.taskDetails.Blur()
				m.help = initializeHelp(stackKeys)
				m.navigationKeys = tableNavigationKeys

			} else if m.taskDetails.Focused() {
				m.stackTable.Blur()
				m.taskTable.Focus()
				m.taskDetails.Blur()
				m.help = initializeHelp(taskKeys)
				m.navigationKeys = tableNavigationKeys

			}
			return m, nil

		case key.Matches(msg, Keys.Right):
			if m.stackTable.Focused() {
				if len(m.stackTable.Rows()) > 0 {
					m.showTasks = true
					m.stackTable.Blur()
					m.taskTable.Focus()
					m.taskDetails.Blur()
					m.help = initializeHelp(taskKeys)
					m.navigationKeys = tableNavigationKeys
					return m, nil
				}
			} else if m.taskTable.Focused() {
				if len(m.taskTable.Rows()) > 0 {
					m.showDetails = true
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = initializeHelp(taskDetailsKeys)
					m.navigationKeys = detailsNavigationKeys
					return m, nil
				}
			} else if m.taskDetails.Focused() {
				m.stackTable.Focus()
				m.taskTable.Blur()
				m.taskDetails.Blur()
				m.help = initializeHelp(stackKeys)
				m.navigationKeys = tableNavigationKeys
				return m, nil
			}

		// Intra-table navigation

		// When we switch to a new stack:
		//		- Empty task box is shown
		//		- Details box is hidden

		// When we switch to a new task:
		//		- Empty details box is shown
		case key.Matches(msg, Keys.Up):
			if m.stackTable.Focused() {
				m.stackTable.MoveUp(1)
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.MoveUp(1)
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.Down):
			if m.stackTable.Focused() {
				m.stackTable.MoveDown(1)
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.MoveDown(1)
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.GotoTop):
			if m.stackTable.Focused() {
				m.stackTable.GotoTop()
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.GotoTop()
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.GotoBottom):
			if m.stackTable.Focused() {
				m.stackTable.GotoBottom()
				m.taskTable.SetCursor(0)
				m.taskDetails.focusIndex = 0
				m.showTasks = false
				m.showDetails = false
				m.updateSelectionData("tasks")
				return m, nil

			} else if m.taskTable.Focused() {
				m.taskTable.GotoBottom()
				m.taskDetails.focusIndex = 0
				m.showDetails = false
				m.updateSelectionData("details")
				return m, nil

			} else if m.taskDetails.Focused() {
				var t tea.Model
				t, cmd = m.taskDetails.Update(msg)
				m.taskDetails = t.(detailsBox)
				return m, cmd
			}

		case key.Matches(msg, Keys.New):
			if m.stackTable.Focused() {
				m.preInputFocus = "stack"
				m.input = initializeInput("stack", entities.Stack{}, 0)

			} else if m.taskTable.Focused() {
				m.preInputFocus = "task"
				newTask := entities.Task{
					StackID: m.data[m.stackTable.Cursor()].ID,
				}
				m.input = initializeInput("task", newTask, 0)

			} else if m.taskDetails.Focused() {
				return m, nil
			}

			m.stackTable.Blur()
			m.taskTable.Blur()
			m.taskDetails.Blur()

			m.showInput = true

			return m, nil

		case key.Matches(msg, Keys.NewRecur):
			if m.taskTable.Focused() {
				m.preInputFocus = "task"
				newTask := entities.Task{
					StackID:            m.data[m.stackTable.Cursor()].ID,
					IsRecurring:        true,
					StartTime:          time.Now(),
					Deadline:           time.Now(),
					RecurrenceInterval: 7,
				}
				m.input = initializeInput("task", newTask, 0)

				m.stackTable.Blur()
				m.taskTable.Blur()
				m.taskDetails.Blur()

				m.showInput = true

				return m, nil
			}

		case key.Matches(msg, Keys.Edit):
			if m.stackTable.Focused() {
				if len(m.stackTable.Rows()) == 0 {
					return m, nil
				}
				m.preInputFocus = "stack"
				m.input = initializeInput("stack", m.data[m.stackTable.Cursor()], 0)
			} else if m.taskTable.Focused() {
				if len(m.taskTable.Rows()) > 0 {
					m.showDetails = true
					m.stackTable.Blur()
					m.taskTable.Blur()
					m.taskDetails.Focus()
					m.help = initializeHelp(taskDetailsKeys)
					m.navigationKeys = detailsNavigationKeys
				}
				return m, nil
			} else if m.taskDetails.Focused() {
				m.preInputFocus = "detail"
				m.input = initializeInput("task", m.data[m.stackTable.Cursor()].Tasks[m.taskTable.Cursor()], m.taskDetails.focusIndex)
			}

			m.stackTable.Blur()
			m.taskTable.Blur()
			m.taskDetails.Blur()

			m.showInput = true

			return m, nil

		//Actual delete operation happens in showDelete conditional at the start of Update() method
		//Here we just trigger the delete confirmation step
		case key.Matches(msg, Keys.Delete):
			if m.stackTable.Focused() {
				m.preInputFocus = "stack"
				m.showDelete = true
				m.deleteConfirmation = initializeDeleteConfirmation()

				return m, nil

			} else if m.taskTable.Focused() {
				stackIndex := m.stackTable.Cursor()

				if len(m.data[stackIndex].Tasks) > 0 {
					m.preInputFocus = "task"
					m.showDelete = true
					m.deleteConfirmation = initializeDeleteConfirmation()

					return m, nil
				}
			}

		case key.Matches(msg, Keys.Toggle):
			//Toggle task finish status
			if m.taskTable.Focused() {
				stackIndex := m.stackTable.Cursor()
				taskIndex := m.taskTable.Cursor()

				var currTask entities.Task
				if len(m.data[stackIndex].Tasks) > 0 {
					stack := m.data[stackIndex]
					currTask = stack.Tasks[taskIndex]

					//For recurring tasks we toggle the status of latest recur task entry
					if currTask.IsRecurring {
						r, count := currTask.LatestRecurTask()
						if count > 0 {
							r.IsFinished = !r.IsFinished
							r.Save()
						}
					} else {
						currTask.IsFinished = !currTask.IsFinished
						currTask.Save()

						if currTask.IsFinished {
							stack.PendingTaskCount--
							stack.Save()
						} else {
							stack.PendingTaskCount++
							stack.Save()
						}

						stack.Tasks[taskIndex] = currTask
						m.data[stackIndex] = stack
					}

					//Changing finish status will lead to reordering, so state has to be preserved
					m.preserveState()
					m.updateSelectionData("stacks")
					return m, nil
				}
			}

		// case key.Matches(msg, Keys.CalendarToggle):
		// 	m.isCalenderView = !m.isCalenderView
		// 	return m, nil

		case key.Matches(msg, Keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		screenWidth = msg.Width
		screenHeight = msg.Height
		if screenHeight < 35 {
			tableViewHeight = 25 - min(35-screenHeight, 10)
		} else {
			tableViewHeight = 25
		}

		//Details box viewport dimensions & section width are set at the time of box creation,
		//after that they have to be manually adjusted
		m.taskDetails.viewport.Width = getDetailsBoxWidth()
		m.taskDetails.viewport.Height = getDetailsBoxHeight()
		m.updateDetailsBoxData(true)

		if m.firstRender {
			//updateSelectionData() is called here instead of being called from Init()
			//since details box rendering requires screen dimensions, which aren't set at the time of Init()
			m.updateSelectionData("stacks")
			m.firstRender = false
		}
	}

	return m, cmd
}

func (m *model) View() string {
	var stackView, taskView, detailView string

	if m.stackTable.Focused() {
		stackView = selectedStackBoxStyle.Render(m.stackView())
		taskView = unselectedBoxStyle.Render(m.taskView())
		detailView = unselectedBoxStyle.Render(m.taskDetails.View())
	} else if m.taskTable.Focused() {
		stackView = unselectedBoxStyle.Render(m.stackView())
		taskView = selectedTaskBoxStyle.Render(m.taskView())
		detailView = unselectedBoxStyle.Render(m.taskDetails.View())
	} else if m.taskDetails.Focused() {
		stackView = unselectedBoxStyle.Render(m.stackView())
		taskView = unselectedBoxStyle.Render(m.taskView())
		detailView = selectedDetailsBoxStyle.Render(m.taskDetails.View())
	} else {
		stackView = unselectedBoxStyle.Render(m.stackView())
		taskView = unselectedBoxStyle.Render(m.taskView())
		detailView = unselectedBoxStyle.Render(m.taskDetails.View())
	}

	// if m.isCalenderView {
	// 	return lipgloss.PlaceHorizontal(screenWidth, lipgloss.Left, initializeCalender(time.Now()).View())
	// }
	viewArr := []string{stackView}
	if m.showTasks {
		viewArr = append(viewArr, taskView)

		if m.showDetails {
			viewArr = append(viewArr, detailView)
		} else if len(m.taskTable.Rows()) > 0 {
			viewArr = append(viewArr, unselectedBoxStyle.Render(getEmptyDetailsView()))
		}
	} else {
		viewArr = append(viewArr, unselectedBoxStyle.Render(getEmptyTaskView()))
	}

	tablesView := lipgloss.JoinHorizontal(lipgloss.Center, viewArr...)

	if m.showDelete {
		return lipgloss.JoinVertical(lipgloss.Left, tablesView, m.deleteConfirmation.View())
	}

	if m.showInput {
		inputFormView := getInputFormStyle().Render(m.input.View())
		tablesView = lipgloss.JoinVertical(lipgloss.Left,
			tablesView,
			inputFormView,
		)
		m.help = initializeHelp(m.input.helpKeys)
	}

	if m.showHelp {
		if !m.showInput {
			navigationHelp := initializeHelp(m.navigationKeys)
			return lipgloss.JoinVertical(lipgloss.Left, tablesView, m.help.View(), navigationHelp.View())
		}
		return lipgloss.JoinVertical(lipgloss.Left, tablesView, m.help.View())
	} else {
		return tablesView
	}
}

func (m model) stackView() string {
	m.stackTable.SetHeight(tableViewHeight)
	return lipgloss.JoinVertical(lipgloss.Center, m.stackTable.View(), m.stackFooter())
}

func (m model) stackFooter() string {
	stackFooterStyle := footerContainerStyle.Copy().
		Width(stackTableWidth)

	info := footerInfoStyle.Render(fmt.Sprintf("%d/%d", m.stackTable.Cursor()+1, len(m.stackTable.Rows())))

	return stackFooterStyle.Render(info)
}

func (m model) taskView() string {
	m.taskTable.SetHeight(tableViewHeight)
	return lipgloss.JoinVertical(lipgloss.Center, m.taskTable.View(), m.taskFooter())
}

func (m model) taskFooter() string {
	taskFooterStyle := footerContainerStyle.Copy().
		Width(taskTableWidth)

	if len(m.taskTable.Rows()) == 0 {
		return taskFooterStyle.Render("Press 'n' to create a new task")
	} else {
		info := footerInfoStyle.Render(fmt.Sprintf("%d/%d", m.taskTable.Cursor()+1, len(m.taskTable.Rows())))
		return taskFooterStyle.Render(info)
	}
}

// Pull new data from database
func (m *model) refreshData() {
	stacks, _ := entities.FetchAllStacks()
	m.data = stacks
	m.updateSelectionData("stacks")
}

// Efficiently update only the required pane
func (m *model) updateSelectionData(category string) {
	var retainIndex bool
	if m.prevState.retainState {
		retainIndex = true
		m.prevState.retainState = false
	}

	switch category {
	case "stacks":
		m.updateStackTableData(retainIndex)
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(true)
	case "tasks":
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(false)
	case "details":
		m.updateDetailsBoxData(false)
	default:
		m.updateStackTableData(retainIndex)
		m.updateTaskTableData(retainIndex)
		m.updateDetailsBoxData(true)
	}
}

func (m *model) updateStackTableData(retainIndex bool) {
	//Set stack view data
	//We pass a slice to stackRows, so the changes (like sorting) that happen there will be reflected in original slice
	m.stackTable.SetRows(stackRows(m.data))

	if retainIndex {
		newIndex := findStackIndex(m.data, m.prevState.stackID)

		if newIndex != -1 {
			m.stackTable.SetCursor(newIndex)
		}
	}
}

func (m *model) updateTaskTableData(retainIndex bool) {
	//Set task view data for selected stack
	stackIndex := m.stackTable.Cursor()
	currStack := m.data[stackIndex]

	//We pass a slice to taskRows, so the changes (like sorting) that happen there will be reflected in original slice
	m.taskTable.SetRows(taskRows(currStack.Tasks))

	if retainIndex {
		newIndex := findTaskIndex(m.data[stackIndex].Tasks, m.prevState.taskID)
		if newIndex != -1 {
			m.taskTable.SetCursor(newIndex)
		}
	}
}

func (m *model) updateDetailsBoxData(preserveOffset bool) {
	stackIndex := m.stackTable.Cursor()
	taskIndex := m.taskTable.Cursor()
	if taskIndex == -1 {
		taskIndex = 0
		m.taskTable.SetCursor(0)
	}

	var currTask entities.Task
	if len(m.data[stackIndex].Tasks) > 0 {
		currTask = m.data[stackIndex].Tasks[taskIndex]
	} else {
		currTask = entities.Task{}
	}

	m.taskDetails.buildDetailsBox(currTask, preserveOffset)
}

// Changing title, deadline, priority or finish status will lead to table reordering
// preserveState() is used to maintain focus on the stack/task that was being edited
func (m *model) preserveState() {
	m.prevState.retainState = true
	stackIndex := m.stackTable.Cursor()
	taskIndex := m.taskTable.Cursor()

	m.prevState.stackID = m.data[m.stackTable.Cursor()].ID
	if len(m.data[stackIndex].Tasks) > 0 {
		m.prevState.taskID = m.data[stackIndex].Tasks[taskIndex].ID
	}
}
