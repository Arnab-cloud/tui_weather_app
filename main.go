package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const DebounceDuration = 500 * time.Millisecond

// var titleStyle = lipgloss.NewStyle().
// 	Bold(true).
// 	Foreground(lipgloss.Color("62")).
// 	MarginBottom(1)

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	MarginTop(1)

var (
	Service *WeatherService
	LogFile *os.File
)

type citySearchResultMsg struct {
	locs []list.Item
}

type weatherSearchResultMsg struct {
	weather *WeatherResponse
}

type debouncedMsg struct {
	id    int
	query string
}

type errorMsg struct {
	err error
}

func (e errorMsg) Error() string { return e.err.Error() }

func (city City) Title() string { return city.Name }
func (city City) Description() string {
	return fmt.Sprintf("%s, Lat: %f, Lon: %f", city.Country, city.Lat, city.Lon)
}
func (city City) FilterValue() string { return city.Name }

type StateModel struct {
	textInput         textinput.Model
	searchResults     list.Model
	curItem           *City
	curWeather        *WeatherResponse
	isFilterOpen      bool
	isFetchingWeather bool
	keys              *itemsKeyMap
	help              help.Model
	err               error
	debounceId        int
	width             int
	height            int
}

func (curM *StateModel) debouncedSearch() tea.Cmd {
	query := curM.textInput.Value()
	curM.debounceId++
	return tea.Tick(DebounceDuration, func(_ time.Time) tea.Msg {
		return debouncedMsg{id: curM.debounceId, query: query}
	})
}

func (curM StateModel) performLocationSearch() tea.Cmd {
	query := curM.textInput.Value()

	return func() tea.Msg {
		if len(query) < 3 {
			return citySearchResultMsg{locs: nil}
		}
		cities, err := Service.ResolveCity(context.Background(), query)
		if err != nil {
			return citySearchResultMsg{locs: nil}
		}
		items := make([]list.Item, len(cities))
		for i, city := range cities {
			items[i] = city
		}
		return citySearchResultMsg{locs: items}
	}
}

func (curM StateModel) performWeatherSearch() tea.Cmd {
	loc := curM.curItem
	if loc == nil {
		return nil
	}
	if loc.Lat == 0 && loc.Lon == 0 {
		return nil
	}
	return func() tea.Msg {
		coord := Coordinates{Lat: curM.curItem.Lat, Lon: curM.curItem.Lon}
		res, err := Service.GetWeather(context.Background(), Location{Name: curM.curItem.Name, Coord: coord, Id: curM.curItem.Id})
		log.Print("called get weather")
		if err != nil {
			log.Printf("error fetching the weather: %s", err)
			return weatherSearchResultMsg{weather: nil}
		}
		return weatherSearchResultMsg{weather: res}
	}
}

func (curM StateModel) Init() tea.Cmd {
	return nil
}

func (curM StateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = nil
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		curM.width, curM.height = msg.Width, msg.Height
		curM.searchResults.SetSize(msg.Width-4, msg.Height-8)

	case citySearchResultMsg:
		curM.searchResults.SetItems(msg.locs)

	case weatherSearchResultMsg:
		curM.curWeather = msg.weather
		curM.isFetchingWeather = false

	case debouncedMsg:
		if curM.debounceId != msg.id {
			return curM, nil
		}
		return curM, curM.performLocationSearch()

	case errorMsg:
		curM.err = msg
		return curM, tea.Quit

	case tea.KeyMsg:

		switch {

		case key.Matches(msg, curM.keys.quitApp) && !curM.textInput.Focused():
			return curM, tea.Quit

		case key.Matches(msg, curM.keys.toggleFilter):
			if !curM.isFilterOpen {
				curM.isFilterOpen = true
				return curM, curM.textInput.Focus()
			}

			if !curM.textInput.Focused() {
				return curM, curM.textInput.Focus()
			}
			curM.isFilterOpen = false
			curM.textInput.Blur()
			return curM, nil

		case key.Matches(msg, curM.keys.choose):
			if i, ok := curM.searchResults.SelectedItem().(City); ok {
				curM.curItem = &i
				curM.isFilterOpen = false
				curM.textInput.Blur()
				curM.isFetchingWeather = true
				cmd = curM.performWeatherSearch()
			}
			return curM, cmd
		case key.Matches(msg, curM.keys.escInput):
			if curM.textInput.Focused() {
				curM.textInput.Blur()
				return curM, nil
			}

			if curM.isFilterOpen {
				curM.isFilterOpen = false
				return curM, nil
			}

			return curM, tea.Quit

		default:
			if !curM.textInput.Focused() {
				curM.searchResults, cmd = curM.searchResults.Update(msg)
				return curM, cmd
			}

			curM.textInput, cmd = curM.textInput.Update(msg)
			return curM, tea.Batch(cmd, curM.debouncedSearch())
		}

	}

	return curM, cmd

}

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

func main() {
	godotenv.Load()

	LogFile, err := tea.LogToFile("logs.log", "debug")
	if err != nil {
		log.Fatalf("Error opening the logfile: %s", err)
	}
	defer LogFile.Close()

	conn, err := sql.Open("sqlite3", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error opening the logfile: %s", err)
	}
	defer conn.Close()

	client := NewWeatherClient(
		os.Getenv("API_KEY"),
		os.Getenv("WEATHER_API"),
		os.Getenv("GEOCODING_API"),
	)

	Service = NewWeatherService(conn, client)

	if _, err := tea.NewProgram(initModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("Error: %s", err)
	}

}

func initModel() StateModel {
	newSearchResults := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	ti := textinput.New()
	ti.Placeholder = "Search for a city"
	ti.CharLimit = 50
	ti.Width = 40

	newModel := StateModel{
		textInput:         ti,
		searchResults:     newSearchResults,
		curItem:           nil,
		curWeather:        nil,
		isFilterOpen:      false,
		isFetchingWeather: false,
		debounceId:        0,
		err:               nil,
		keys:              newItemsKeyMap(),
		help:              help.New(),
	}
	newModel.searchResults.Title = "Find Cities"
	newModel.searchResults.SetShowFilter(false)
	newModel.searchResults.SetFilteringEnabled(false)

	return newModel
}

func renderWeather(weather *WeatherResponse, width, height int) string {
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

func getWeatherEmoji(icon string) string {
	if emoji, ok := emojiMap[icon]; ok {
		return emoji
	}
	return "🌤️"
}
