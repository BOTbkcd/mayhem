package tui

import (
	"github.com/BOTbkcd/mayhem/entities"
	"github.com/BOTbkcd/mayhem/utils"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type trelloSync struct {
	data           entities.SyncInfo
	trelloBoardURL tea.Model
	trelloKey      tea.Model
	trelloToken    tea.Model
	focusIndex     int
}

var trelloSyncKeys = keyMap{
	Toggle: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("'tab'", "next field"),
	),
	ReverseToggle: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("'shift+tab'", "previous field"),
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

func initializeTrelloSync(syncData entities.SyncInfo) tea.Model {
	t := trelloSync{
		data:           syncData,
		trelloBoardURL: initializeTextInput(syncData.BoardURL, "Enter Trello Board URL", 0, 0, goToSyncWithVal),
		trelloKey:      initializeTextInput(syncData.Key, "Enter Trello API Key", 0, 0, goToSyncWithVal),
		trelloToken:    initializeTextInput(syncData.Token, "Enter Trello API Token", 0, 30, goToSyncWithVal),
	}

	return t
}

func (m trelloSync) Init() tea.Cmd {
	return nil
}

func (m trelloSync) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.ReverseToggle):
			if m.focusIndex > 0 {
				m.focusIndex--
			}
			return m, nil

		case key.Matches(msg, Keys.Toggle):
			if m.focusIndex < 4 {
				m.focusIndex++
			}
			return m, nil

		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, Keys.Enter):
			switch m.focusIndex {
			case 0:
				utils.SyncBoardData(m.data.BoardURL, m.data.Key, m.data.Token)
			case 1:
				board := utils.GenerateNewBoard(m.data.Key, m.data.Token)
				m.data.BoardURL = board.Url
				m.data = m.data.Save()
			case 2:
				m.trelloBoardURL, cmd = m.trelloBoardURL.Update(msg)
			case 3:
				m.trelloKey, cmd = m.trelloKey.Update(msg)
			case 4:
				m.trelloToken, cmd = m.trelloToken.Update(msg)
			}

		default:
			switch m.focusIndex {
			case 2:
				m.trelloBoardURL, cmd = m.trelloBoardURL.Update(msg)
			case 3:
				m.trelloKey, cmd = m.trelloKey.Update(msg)
			case 4:
				m.trelloToken, cmd = m.trelloToken.Update(msg)
			}
		}

	case goToSyncMsg:
		switch m.focusIndex {
		case 2:
			m.data.BoardURL = msg.value.(string)
		case 3:
			m.data.Key = msg.value.(string)
		case 4:
			m.data.Token = msg.value.(string)
		}

		m.data = m.data.Save()

	default:
		switch m.focusIndex {
		case 2:
			m.trelloBoardURL, cmd = m.trelloBoardURL.Update(msg)
		case 3:
			m.trelloKey, cmd = m.trelloKey.Update(msg)
		case 4:
			m.trelloToken, cmd = m.trelloToken.Update(msg)
		}
	}

	return m, cmd
}

func (m trelloSync) View() string {
	// blurredButton := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	// focusedButton := focusedStyle.Copy().Render("[ Submit ]")

	// var button *string
	// if m.focusIndex == len(m.fieldMap) {
	// 	button = &focusedButton
	// } else {
	// 	button = &blurredButton
	// }

	// fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	syncButtonStyle := unselectedButtonStyle
	generateButtonStyle := unselectedButtonStyle

	fieldStyle := inputFieldStyle.Copy().Width(40)
	urlField := fieldStyle.Render(m.data.BoardURL)
	keyField := fieldStyle.Render(m.data.Key)
	tokenField := fieldStyle.Render(m.data.Token)

	switch m.focusIndex {
	case 0:
		syncButtonStyle = selectedButtonStyle
	case 1:
		generateButtonStyle = selectedButtonStyle
	case 2:
		urlField = fieldStyle.Copy().BorderForeground(timeFocusColor).Render(m.trelloBoardURL.View())
	case 3:
		keyField = fieldStyle.Copy().BorderForeground(timeFocusColor).Render(m.trelloKey.View())
	case 4:
		tokenField = fieldStyle.Copy().BorderForeground(timeFocusColor).Render(m.trelloToken.View())
	}

	syncButton := syncButtonStyle.Render("Sync")
	generateButton := generateButtonStyle.Render("Generate")

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, syncButton, generateButton)

	boardUrl := lipgloss.JoinHorizontal(lipgloss.Center, "Trello Board Url: ", urlField)
	key := lipgloss.JoinHorizontal(lipgloss.Center, "  Trello API Key: ", keyField)
	token := lipgloss.JoinHorizontal(lipgloss.Center, "Trello API Token: ", tokenField)

	fields := lipgloss.JoinVertical(lipgloss.Left, boardUrl, key, token)

	return lipgloss.JoinVertical(lipgloss.Center, buttons, fields)
}
