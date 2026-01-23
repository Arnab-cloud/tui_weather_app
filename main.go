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

const DebounceDuration = time.Second

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

type errorMsg struct {
	err error
}

func (e errorMsg) Error() string { return e.err.Error() }

type item struct {
	title, desc string
}

func citiesToItems(locs []City) []list.Item {
	cities := make([]list.Item, len(locs))
	for i, loc := range locs {
		cities[i] = item{title: loc.Name, desc: fmt.Sprintf("%s, Lat: %f, Lon: %f", loc.Country, loc.Lat, loc.Lon)}
	}
	return cities
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type StateModel struct {
	textInput     textinput.Model
	searchResults list.Model
	curitem       *item
	isFilterOpen  bool
	keys          *itemsKeyMap
	help          help.Model
	err           error
	debounceId    int
	width         int
	height        int
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

		return searchMsgResult{locs: citiesToItems(cities)}
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
			if i, ok := curM.searchResults.SelectedItem().(item); ok {
				curM.curitem = &i
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
			return curM, tea.Batch(cmd, curM.performSearch())
		}

	}

	return curM, cmd

}

func (curM StateModel) View() string {

	if !curM.isFilterOpen && curM.curitem != nil {
		return fmt.Sprintf("Selected: %s, %s\nFetching weather...\n",
			curM.curitem.title, curM.curitem.desc)
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
	newItems := []list.Item{
		// item{title: "Coffee", desc: "Freshly brewed hot drink"},
		// item{title: "Pizza", desc: "Cheesy slice with pepperoni"},
		// item{title: "Sushi", desc: "Rice rolls with fresh fish"},
		// item{title: "Burger", desc: "Grilled beef patty with toppings"},
		// item{title: "Pasta", desc: "Italian noodles in tomato sauce"},

		// item{title: "Library", desc: "Quiet place full of books"},
		// item{title: "Beach", desc: "Sandy shore by the ocean"},
		// item{title: "Mountain", desc: "High elevation hiking spot"},
		// item{title: "Cafe", desc: "Small shop for coffee and snacks"},
		// item{title: "Park", desc: "Green space for relaxing walks"},

		// item{title: "Laptop", desc: "Portable computer for work"},
		// item{title: "Headphones", desc: "Noise-canceling audio gear"},
		// item{title: "Backpack", desc: "Bag for carrying essentials"},
		// item{title: "Smartphone", desc: "Touchscreen mobile device"},
		// item{title: "Notebook", desc: "Paper pad for writing notes"},

		// item{title: "Gym", desc: "Place to exercise and train"},
		// item{title: "Museum", desc: "Exhibits of history and art"},
		// item{title: "Airport", desc: "Hub for air travel"},
		// item{title: "Restaurant", desc: "Dine-in food service"},
		// item{title: "Cinema", desc: "Theater for watching movies"},
	}

	newsearchResults := list.New(newItems, list.NewDefaultDelegate(), 0, 0)
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
