package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
		return errorStyle.Render(fmt.Sprintf("❌ Error: %v", curM.err))
	}

	if curM.isFilterOpen {
		var s strings.Builder
		s.WriteString(titleStyle.Render("🌤️  Weather Search"))
		s.WriteString("\n\n")
		s.WriteString(curM.textInput.View())
		s.WriteString("\n\n")
		s.WriteString(curM.searchResults.View())
		s.WriteString("\n")

		return s.String()
	}

	if curM.curItem != nil && curM.isFetchingWeather {
		return fmt.Sprintf("Selected: %s, %s\nFetching weather...\n",
			curM.curItem.Title(), curM.curItem.Description())
	}

	if curM.curItem != nil && curM.curWeather != nil {
		return renderWeather(curM.curWeather)
	}

	return "Nothing to show here"

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

func renderWeather(weather *WeatherResponse) string {
	var b strings.Builder
	caser := cases.Title(language.English)

	// Header with city name and country
	header := fmt.Sprintf("📍 %s, %s", weather.Name, weather.Sys.Country)
	b.WriteString(cityStyle.Render(header))
	b.WriteString("\n\n")

	// Main temperature display
	temp := fmt.Sprintf("%.1f°C", weather.Main.Temp)
	b.WriteString(tempStyle.Render(temp))
	b.WriteString("\n")

	// Weather description with icon
	if len(weather.Weather) > 0 {
		weatherDesc := getWeatherEmoji(weather.Weather[0].Icon) + " " +
			caser.String(weather.Weather[0].Desc)
		b.WriteString(descStyle.Render(weatherDesc))
		b.WriteString("\n\n")
	}

	// Feels like temperature
	feelsLike := fmt.Sprintf("%s%s",
		labelStyle.Render("Feels like:"),
		valueStyle.Render(fmt.Sprintf("%.1f°C", weather.Main.FeelsLike)))
	b.WriteString(feelsLike)
	b.WriteString("\n")

	// Temperature range
	tempRange := fmt.Sprintf("%s%s",
		labelStyle.Render("Range:"),
		valueStyle.Render(fmt.Sprintf("%.1f°C - %.1f°C",
			weather.Main.TempMin, weather.Main.TempMax)))
	b.WriteString(tempRange)
	b.WriteString("\n\n")

	// Weather details box
	var details strings.Builder

	details.WriteString(formatDetail("💧 Humidity", fmt.Sprintf("%d%%", weather.Main.Humidity)))
	details.WriteString(formatDetail("🌬  Wind Speed", fmt.Sprintf("%.1f m/s", weather.Wind.Speed)))
	details.WriteString(formatDetail("🧭 Wind Direction", getWindDirection(weather.Wind.Deg)))
	details.WriteString(formatDetail("🔽 Pressure", fmt.Sprintf("%d hPa", weather.Main.Pressure)))
	details.WriteString(formatDetail("👁  Visibility", fmt.Sprintf("%d m", weather.Vis)))
	if weather.Clouds > 0 {
		details.WriteString(formatDetail("☁️  Clouds", fmt.Sprintf("%d%%", weather.Clouds)))
	}

	b.WriteString(boxStyle.Render(details.String()))
	b.WriteString("\n")

	// Sun times
	sunrise := time.Unix(weather.Sys.Sunrise, 0).Format("15:04")
	sunset := time.Unix(weather.Sys.Sunset, 0).Format("15:04")

	sunInfo := fmt.Sprintf("🌅 Sunrise: %s  |  🌇 Sunset: %s",
		valueStyle.Render(sunrise),
		valueStyle.Render(sunset))
	b.WriteString(sunInfo)
	b.WriteString("\n\n")

	// Footer
	footer := mutedColorStyle.Render("Press 'q' to quit | Press 's' to search new city")
	b.WriteString(footer)

	return b.String()
}

func getWeatherEmoji(icon string) string {
	if emoji, ok := emojiMap[icon]; ok {
		return emoji
	}
	return "🌤️"
}

func formatDetail(label, value string) string {
	return fmt.Sprintf("%s%s\n",
		labelStyle.Render(label+":"),
		valueStyle.Render(value))
}

func getWindDirection(deg int) string {
	directions := []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
		"S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}
	index := int((float64(deg) + 11.25) / 22.5)
	return directions[index%16] + fmt.Sprintf(" (%d°)", deg)
}
