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
	help         key.Binding
}

func newItemsKeyMap() *itemsKeyMap {
	return &itemsKeyMap{
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
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
			key.WithHelp("k/↑", "up"),
		),
		down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "down"),
		),
		back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k itemsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.toggleFilter, k.quit}
}

func (k itemsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down, k.choose},
		{k.toggleFilter, k.back, k.quit},
	}
}

func (k itemsKeyMap) GetContextualHelp(isFilterOpen, isInputFocused bool) []key.Binding {
	if !isFilterOpen {
		return []key.Binding{k.toggleFilter, k.quit}
	}

	if isInputFocused {
		escBlur := key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "blur input"),
		)
		return []key.Binding{k.toggleFilter, escBlur}
	}

	escExit := key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "exit filter"),
	)
	return []key.Binding{k.up, k.down, k.choose, k.toggleFilter, escExit, k.quit}
}

type contextualKeyMap struct {
	bindings []key.Binding
}

func (k *contextualKeyMap) ShortHelp() []key.Binding {
	return k.bindings
}

func (k *contextualKeyMap) FullHelp() [][]key.Binding {
	// Split bindings into rows of max 4 items for better layout
	if len(k.bindings) <= 4 {
		return [][]key.Binding{k.bindings}
	}

	mid := (len(k.bindings) + 1) / 2
	return [][]key.Binding{
		k.bindings[:mid],
		k.bindings[mid:],
	}
}
