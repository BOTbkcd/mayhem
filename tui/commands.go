package tui

import (
	"github.com/BOTbkcd/mayhem/entities"
	tea "github.com/charmbracelet/bubbletea"
)

type goToMainMsg struct {
	value interface{}
}

func goToMainCmd() tea.Msg {
	return goToMainMsg{
		value: "",
	}
}

func goToMainWithVal(value interface{}) tea.Cmd {
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

type goToSyncMsg struct {
	value interface{}
}

func goToSyncWithVal(value interface{}) tea.Cmd {
	return func() tea.Msg {
		return goToSyncMsg{value: value}
	}
}

/*
************* Database Operations *****************
 */
type updatedData []entities.Stack

func fetchUpdatedStacks() tea.Msg {
	return updatedData(entities.FetchAllStacks())
}

type syncInfo entities.SyncInfo

func fetchSyncInfo() tea.Msg {
	return syncInfo(entities.FetchSyncInfo()[0])
}

/*
************* Trello API Calls *****************
 */

func fetchTrelloLists(syncInfo entities.SyncInfo) tea.Cmd {
	return func() tea.Msg {
		return ""
	}
}
