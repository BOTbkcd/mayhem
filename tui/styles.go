package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// func FontColor(str, color string) string {
// 	var term = termenv.ColorProfile()
// 	return termenv.String(str).Foreground(term.Color(color)).Bold().String()
// }

var (
	screenWidth     int
	screenHeight    int
	tableViewHeight = 25
	stackTableWidth = 26 //22: column width + 2*2: column padding
	taskTableWidth  = 67 //59: column widths + 2*4: column paddings
	dash            = "–"
)

var (
	stackBorderColor     = lipgloss.Color("#019187")
	taskBorderColor      = lipgloss.Color("#f1b44c")
	detailsBorderColor   = lipgloss.Color("#6192bc")
	inputFormBorderColor = lipgloss.Color("#325b84")

	stackSelectionColor   = lipgloss.Color("#019187")
	taskSelectionColor    = lipgloss.Color("#f1b44c")
	detailsSelectionColor = lipgloss.Color("#333c4d")

	highlightedBackgroundColor = lipgloss.Color("#f97171")
	highlightedTextColor       = lipgloss.Color("#4e4e4e")
	inputFormColor             = lipgloss.Color("#5ac7c7")
	timeFocusColor             = lipgloss.Color("#FFFF00")
	unfocusedColor             = lipgloss.Color("#898989")
	whiteColor                 = lipgloss.Color("#ffffff")
)
var (
	selectedBoxStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.ThickBorder())

	selectedStackBoxStyle = selectedBoxStyle.Copy().
				BorderForeground(stackBorderColor)

	selectedTaskBoxStyle = selectedBoxStyle.Copy().
				BorderForeground(taskBorderColor)

	selectedDetailsBoxStyle = selectedBoxStyle.Copy().
				BorderForeground(detailsBorderColor)

	unselectedBoxStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(unfocusedColor)

	columnHeaderStyle = table.DefaultStyles().Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(unfocusedColor).
				BorderBottom(true).
				Bold(true)

	stackSelectedRowStyle = table.DefaultStyles().Selected.
				Foreground(highlightedTextColor).
				Background(stackSelectionColor).
				Bold(false)

	taskSelectedRowStyle = stackSelectedRowStyle.Copy().
				Background(taskSelectionColor)

	footerInfoStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color("#1c2c4c"))

	footerContainerStyle = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Background(lipgloss.Color("#3e424b"))

	highlightedTextStyle = lipgloss.NewStyle().
				Bold(true).
				Italic(true).
				Foreground(highlightedTextColor).
				Background(highlightedBackgroundColor).
				Padding(0, 1).
				MarginTop(1)

	inputFormStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(inputFormBorderColor).
			Padding(0, 1)

	textInputStyle   = lipgloss.NewStyle().Foreground(inputFormColor)
	placeHolderStyle = lipgloss.NewStyle().Foreground(unfocusedColor)
)

// Since width is dynamic, we have to append it to the style before usage

func getInputFormStyle() lipgloss.Style {
	//Subtract 2 for padding on each side
	return inputFormStyle.Width(screenWidth - 2)
}

func getTableStyle(tableType string) table.Styles {
	s := table.DefaultStyles()
	s.Header = columnHeaderStyle

	switch tableType {
	case "stack":
		s.Selected = stackSelectedRowStyle
	case "task":
		s.Selected = taskSelectedRowStyle
	}

	return s
}

func getEmptyTaskStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Width(60).
		Height(tableViewHeight + 3) //3 is added to account for header & footer height
}

func getEmptyDetailsStyle() lipgloss.Style {
	return getDetailsBoxStyle().
		Height(tableViewHeight + 3).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
}

func getDetailsBoxWidth() int {
	return screenWidth - (stackTableWidth + taskTableWidth + 3*2) //each of the 3 boxes have left & right borders
}
func getDetailsBoxHeight() int {
	return tableViewHeight
}

func getDetailsBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(getDetailsBoxWidth()).
		Height(tableViewHeight + 2)
}

func getDetailsItemStyle(isSelected bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Padding(0, 0, 1, 0).
		Width(getDetailsBoxWidth() - 2)

	if isSelected {
		style.Background(detailsSelectionColor)
	}

	return style
}

// Applying padding (0,1) to detail items causes issue with description text alignment
// To avoid that an additional container is used for detail items
func getItemContainerStyle(isSelected bool) lipgloss.Style {
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Width(getDetailsBoxWidth())

	if isSelected {
		style.Background(detailsSelectionColor)
	}

	return style
}
