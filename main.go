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
)

const DebounceDuration = 500 * time.Millisecond

var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("62")).
	MarginBottom(1)

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("241")).
	MarginTop(1)

var (
	Service *WeatherService
	LogFile *os.File
)

type searchMsgResult struct {
	locs []list.Item
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
	textInput     textinput.Model
	searchResults list.Model
	curItem       *City
	isFilterOpen  bool
	keys          *itemsKeyMap
	help          help.Model
	err           error
	debounceId    int
}

func (curM *StateModel) debouncedSearch() tea.Cmd {
	query := curM.textInput.Value()
	curM.debounceId++
	return tea.Tick(DebounceDuration, func(_ time.Time) tea.Msg {
		return debouncedMsg{id: curM.debounceId, query: query}
	})
}

func (curM StateModel) performSearch() tea.Cmd {
	query := curM.textInput.Value()

	return func() tea.Msg {
		if len(query) < 3 {
			return searchMsgResult{locs: nil}
		}
		cities, err := Service.ResolveCity(context.Background(), query)
		if err != nil {
			return searchMsgResult{locs: nil}
		}
		items := make([]list.Item, len(cities))
		for i, city := range cities {
			items[i] = city
		}
		return searchMsgResult{locs: items}
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

	case searchMsgResult:
		curM.searchResults.SetItems(msg.locs)

	case debouncedMsg:
		if curM.debounceId != msg.id {
			return curM, nil
		}
		return curM, curM.performSearch()

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
			}
			return curM, nil
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

	if !curM.isFilterOpen && curM.curItem != nil {
		return fmt.Sprintf("Selected: %s, %s\nFetching weather...\n",
			curM.curItem.Title(), curM.curItem.Description())
	}

	if !curM.isFilterOpen {
		return "Nothing to show here"
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render("🌤️  Weather Search"))
	s.WriteString("\n\n")
	s.WriteString(curM.textInput.View())
	s.WriteString("\n\n")
	s.WriteString(curM.searchResults.View())
	s.WriteString("\n")

	return s.String()
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
	newsearchResults := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	ti := textinput.New()
	ti.Placeholder = "Search for a city"
	ti.CharLimit = 50
	ti.Width = 40

	newModel := StateModel{
		textInput:     ti,
		searchResults: newsearchResults,
		isFilterOpen:  false,
		keys:          newItemsKeyMap(),
		help:          help.New(),
	}
	newModel.searchResults.Title = "Find Cities"
	newModel.searchResults.SetShowFilter(false)
	newModel.searchResults.SetFilteringEnabled(false)

	return newModel
}
