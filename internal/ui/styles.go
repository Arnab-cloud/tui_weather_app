package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Primary colors
	fg          = lipgloss.Color("#a9b1d6")
	cyan        = lipgloss.Color("#7dcfff")
	magenta     = lipgloss.Color("#bb9af7")
	yellow      = lipgloss.Color("#e0af68")
	green       = lipgloss.Color("#9ece6a")
	blue        = lipgloss.Color("#7aa2f7")
	borderColor = lipgloss.Color("#414868")
	comment     = lipgloss.Color("#565f89")
	errorColor  = lipgloss.Color("#E06C75")
	white       = lipgloss.Color("#FFFFFF")
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true).
			Padding(0, 1)

	windowStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center)

	hiLoLabelStyle = lipgloss.NewStyle().
			Foreground(comment).
			MarginRight(1)

	hiLoValueStyle = lipgloss.NewStyle().
			Foreground(fg).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

var emojiMap = map[string]string{
	"01d": "☀️", // clear sky day
	"01n": "🌙",  // clear sky night
	"02d": "⛅",  // few clouds day
	"02n": "☁️", // few clouds night
	"03d": "☁️", // scattered clouds
	"03n": "☁️",
	"04d": "☁️", // broken clouds
	"04n": "☁️",
	"09d": "🌧️", // shower rain
	"09n": "🌧️",
	"10d": "🌦️", // rain
	"10n": "🌧️",
	"11d": "⛈️", // thunderstorm
	"11n": "⛈️",
	"13d": "❄️", // snow
	"13n": "❄️",
	"50d": "🌫️", // mist
	"50n": "🌫️",
}

func getWeatherEmoji(icon string) string {
	if emoji, ok := emojiMap[icon]; ok {
		return emoji
	}
	return "🌤️"
}
