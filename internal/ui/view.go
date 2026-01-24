package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (curM StateModel) View() string {

	if curM.err != nil {
		return windowStyle.
			Width(curM.width).
			Height(curM.height).
			Render(errorStyle.Render(fmt.Sprintf("❌ Error: %v", curM.err)))
	}

	if curM.isFilterOpen {
		searchContent := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("🌤️ Weather Search"),
			curM.textInput.View(),
			"",
			curM.searchResults.View(),
		)
		return windowStyle.
			Width(curM.width).
			Height(curM.height).
			Render(searchContent)
	}

	if curM.curItem != nil && curM.curWeather != nil {
		return renderWeather(curM.curWeather, curM.width, curM.height)
	}

	return windowStyle.
		Width(curM.width).
		Height(curM.height).
		Render("Loading...")

}
