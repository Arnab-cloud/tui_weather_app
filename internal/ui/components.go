package ui

import (
	"fmt"
	"github/Arnab-cloud/tui_weather_app/internal/weather"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

func renderWeather(weather *weather.WeatherResponse, width, height int) string {
	// Hero section with location and temperature
	fig := figure.NewColorFigure(fmt.Sprintf("%.1f", weather.Main.Temp), "slant", "yellow", true)
	bigTemp := fig.String()

	locationStyle := lipgloss.NewStyle().Foreground(white).Bold(true)
	location := locationStyle.Render(fmt.Sprintf("📍 %s, %s", weather.Name, weather.Sys.Country))

	weatherDesc := lipgloss.JoinHorizontal(lipgloss.Center,
		getWeatherEmoji(weather.Weather[0].Icon),
		lipgloss.NewStyle().MarginLeft(2).Foreground(fg).Render(weather.Weather[0].Desc),
	)

	heroLeft := lipgloss.JoinVertical(lipgloss.Left,
		location,
		weatherDesc,
		"",
		formatHiLo(weather.Main.TempMax, weather.Main.TempMin),
	)

	hero := renderSection("",
		lipgloss.JoinHorizontal(lipgloss.Center,
			heroLeft,
			lipgloss.NewStyle().Width(10).Render(""),
			lipgloss.NewStyle().Foreground(yellow).Render(bigTemp),
		),
		width-10,
		yellow,
	)

	// Atmosphere grid
	colWidth := (width / 3) - 6
	atmRow1 := lipgloss.JoinHorizontal(lipgloss.Top,
		renderDataPoint("🌡️ Feels Like", fmt.Sprintf("%.1f°C", weather.Main.FeelsLike), colWidth),
		renderDataPoint("💧 Humidity", fmt.Sprintf("%d%%", weather.Main.Humidity), colWidth),
		renderDataPoint("🌬️ Wind", fmt.Sprintf("%.1f m/s", weather.Wind.Speed), colWidth),
	)

	atmRow2 := lipgloss.JoinHorizontal(lipgloss.Top,
		renderDataPoint("⏲️ Pressure", fmt.Sprintf("%d hPa", weather.Main.Pressure), colWidth),
		renderDataPoint("👁️ Visibility", fmt.Sprintf("%.1f km", float64(weather.Vis)/1000), colWidth),
		renderDataPoint("☁️ Cloudiness", fmt.Sprintf("%d%%", weather.Clouds), colWidth),
	)

	atmosphere := renderSection("Atmosphere",
		lipgloss.JoinVertical(lipgloss.Left, atmRow1, atmRow2),
		width-10,
		cyan,
	)

	// Sun times
	halfWidth := (width / 2) - 8
	sunContent := lipgloss.JoinHorizontal(lipgloss.Top,
		renderDataPoint("🌅 Sunrise", time.Unix(weather.Sys.Sunrise, 0).Format("03:04 PM"), halfWidth/2),
		renderDataPoint("🌇 Sunset", time.Unix(weather.Sys.Sunset, 0).Format("03:04 PM"), halfWidth/2),
	)

	sunSection := renderSection("Sun Times", sunContent, width-10, blue)

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(comment).
		Render("Press 'q' to quit • 's' to search")

	// Assemble full view
	fullView := lipgloss.JoinVertical(lipgloss.Left,
		hero,
		"",
		atmosphere,
		"",
		sunSection,
		"",
		footer,
	)

	return windowStyle.
		Width(width).
		Height(height).
		Render(fullView)
}

func renderSection(title, content string, width int, color lipgloss.Color) string {
	sectionTitle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Padding(0, 1).
		Render(title)

	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Render(content)

	if title == "" {
		return border
	}

	// We render the title, then "move" it slightly down or just place it above
	return lipgloss.JoinVertical(
		lipgloss.Left,
		sectionTitle,
		border,
	)
}

func renderDataPoint(label, value string, width int) string {
	l := lipgloss.NewStyle().Foreground(magenta).Render(label)
	v := lipgloss.NewStyle().Foreground(green).Bold(true).Render(value)

	return lipgloss.NewStyle().
		Width(width).
		Padding(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, l, v))
}

func formatHiLo(hi, lo float64) string {
	high := lipgloss.JoinHorizontal(lipgloss.Left,
		hiLoLabelStyle.Render("H:"),
		hiLoValueStyle.Render(fmt.Sprintf("%.0f°", hi)),
	)

	low := lipgloss.JoinHorizontal(lipgloss.Left,
		hiLoLabelStyle.Render("L:"),
		hiLoValueStyle.Render(fmt.Sprintf("%.0f°", lo)),
	)

	return lipgloss.JoinHorizontal(lipgloss.Left, high, "  ", low)
}
