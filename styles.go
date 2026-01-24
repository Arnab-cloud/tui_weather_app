package main

import "github.com/charmbracelet/lipgloss"

var (
	// Color scheme
	primaryColor   = lipgloss.Color("#61AFEF")
	secondaryColor = lipgloss.Color("#98C379")
	accentColor    = lipgloss.Color("#E5C07B")
	errorColor     = lipgloss.Color("#E06C75")
	mutedColor     = lipgloss.Color("#5C6370")
	bgColor        = lipgloss.Color("#282C34")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	cityStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)
		// FontSize(24)

	tempStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)
		// FontSize(36)

	labelStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(15)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ABB2BF")).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1)

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C678DD")).
			Italic(true)

	mutedColorStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor)
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
