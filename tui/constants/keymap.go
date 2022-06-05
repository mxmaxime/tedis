package constants

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	Create key.Binding
	Enter  key.Binding
	Rename key.Binding
	Back   key.Binding
	Save   key.Binding
	Cancel key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Create: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "insert"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Save: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "write/save"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "cancel"),
	),
}
