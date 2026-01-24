package main

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

// A "Card" to hold the content
// cardStyle = lipgloss.NewStyle().
// Border(lipgloss.RoundedBorder()).
// BorderForeground(primaryColor).
// Padding(1, 4).
// Width(60)

// Large Temperature Display
// bigTempStyle = lipgloss.NewStyle().
// Foreground(accentColor).
// Bold(true).
// Padding(0, 1).
// MarginRight(4).
// SetString("") // Placeholder

// Sidebar/Details column
// columnStyle = lipgloss.NewStyle().
// 		Width(25)

// Use Padding and Width to create structural "size"
// heroColumn = lipgloss.NewStyle().
// 		Padding(1, 2).
// 		Border(lipgloss.NormalBorder(), false, true, false, false). // Right border only
// 		BorderForeground(mutedColor)

// Info rows
// infoLabel = lipgloss.NewStyle().Foreground(mutedColor)
// infoValue = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

// cityStyle = lipgloss.NewStyle().
// 		Foreground(secondaryColor).
// 		Bold(true)

// tempStyle = lipgloss.NewStyle().
// Foreground(accentColor).
// Bold(true)

// labelStyle = lipgloss.NewStyle().
// Foreground(mutedColor).
// Width(15)

// valueStyle = lipgloss.NewStyle().
// 		Foreground(lipgloss.Color("#ABB2BF")).
// 		Bold(true)

// boxStyle = lipgloss.NewStyle().
// 		Border(lipgloss.RoundedBorder()).
// 		BorderForeground(primaryColor).
// 		Padding(1, 2).
// 		MarginTop(1)

// descStyle = lipgloss.NewStyle().
// Foreground(lipgloss.Color("#C678DD")).
// Italic(true)

// mutedColorStyle = lipgloss.NewStyle().
// 		Foreground(mutedCo
