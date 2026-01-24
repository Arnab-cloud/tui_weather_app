package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type itemsKeyMap struct {
	choose       key.Binding
	toggleFilter key.Binding // "/"
	up           key.Binding // "k"
	down         key.Binding // "j"
	back         key.Binding // "esc"
	quit         key.Binding // "q"
}

func newItemsKeyMap() *itemsKeyMap {
	return &itemsKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		toggleFilter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j", "down"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back/cancel"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// Help helpers
func (k itemsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.toggleFilter, k.back, k.choose, k.quit}
}

func (k itemsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down, k.choose},
		{k.toggleFilter, k.back, k.quit},
	}
}
