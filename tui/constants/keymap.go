package constants

import "github.com/charmbracelet/bubbles/key"

type keymap struct {
	Create key.Binding
	Enter  key.Binding
	Rename key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}
