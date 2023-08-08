package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	CalendarToggle key.Binding
	Up             key.Binding
	Down           key.Binding
	GotoTop        key.Binding
	GotoBottom     key.Binding
	Left           key.Binding
	Right          key.Binding
	New            key.Binding
	NewRecur       key.Binding
	Edit           key.Binding
	Move           key.Binding
	Enter          key.Binding
	Save           key.Binding
	Toggle         key.Binding
	ReverseToggle  key.Binding
	Delete         key.Binding
	Return         key.Binding
	Help           key.Binding
	Quit           key.Binding
	Exit           key.Binding
}

var Keys = keyMap{
	CalendarToggle: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("'c'", "calendar view"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("'↑/k'", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("'↓/j'", "move down"),
	),
	GotoTop: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("'g'", "go to top"),
	),
	GotoBottom: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("'G'", "go to bottom"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("'←/h'", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("'→/l'", "move right"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("'n'", "new"),
	),
	NewRecur: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("'r'", "new recurring"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("'e'", "edit"),
	),
	Move: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("'m'", "move"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("'enter'", "enter"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("'tab'", "toggle"),
	),
	ReverseToggle: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("'shift+tab'", "toggle"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("'x'", "delete 🗑"),
	),
	Return: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("'esc'", "return"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("'?'", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("'q'", "quit"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("'ctrl+c'", "exit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.CalendarToggle,
		k.Toggle,
		k.ReverseToggle,
		k.New,
		k.NewRecur,
		k.Edit,
		k.Enter,
		k.Save,
		k.Delete,
		k.Move,
		k.Return,
		k.Up,
		k.Down,
		k.GotoTop,
		k.GotoBottom,
		k.Left,
		k.Right,
		k.Help,
		k.Quit,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
