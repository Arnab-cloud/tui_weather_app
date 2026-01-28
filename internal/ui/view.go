package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (curM StateModel) View() string {
	var content string

	helpView := curM.renderContextualHelp()
	height := max(curM.height-lipgloss.Height(helpView), 0)

	if curM.err != nil {
		content = windowStyle.
			Width(curM.width).
			Height(height).
			Render(errorStyle.Render(fmt.Sprintf("❌ Error: %v", curM.err)))
	} else if curM.isFilterOpen {
		searchContent := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("🌤️ Weather Search"),
			curM.textInput.View(),
			"",
			curM.searchResults.View(),
		)
		content = windowStyle.
			Width(curM.width).
			Height(height).
			Render(searchContent)
	} else if curM.curItem != nil && curM.curWeather != nil {
		content = renderWeather(curM.curWeather, curM.width, height)
	} else {
		content = windowStyle.
			Width(curM.width).
			Height(height).
			Render("Loading...")
	}

	return lipgloss.JoinVertical(lipgloss.Left, content, helpView)
}

func (curM StateModel) renderContextualHelp() string {
	contextualBindings := curM.keys.GetContextualHelp(curM.isFilterOpen, curM.textInput.Focused())

	helpKeys := &contextualKeyMap{bindings: contextualBindings}

	return curM.help.View(helpKeys)
}
