package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type itemsKeyMap struct {
	choose       key.Binding
	toggleFilter key.Binding
	quitApp      key.Binding
	escInput     key.Binding
}

func (dKeyMp itemsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{dKeyMp.escInput, dKeyMp.choose, dKeyMp.toggleFilter, dKeyMp.quitApp}
}

func (dKeyMp itemsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{dKeyMp.choose, dKeyMp.toggleFilter, dKeyMp.quitApp, dKeyMp.escInput},
	}
}

func newItemsKeyMap() *itemsKeyMap {
	return &itemsKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "Enter to select"),
		),
		toggleFilter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "Toggle filter"),
		),
		quitApp: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "close app"),
		),
		escInput: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "esc from input"),
		),
	}
}
