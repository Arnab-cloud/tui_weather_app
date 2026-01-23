package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const DebounceDuration = time.Second

var appStyle = lipgloss.NewStyle().Padding(1, 2)

var (
	Service *WeatherService
	LogFile *os.File
)

type errorMsg struct {
	err error
}

func (e errorMsg) Error() string { return e.err.Error() }

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type StateModel struct {
	locationList list.Model
	curitem      *item
	isFilterOpen bool
	keys         *itemsKeyMap
	help         help.Model
	err          error
	debounceId   int
}

func (curM StateModel) Init() tea.Cmd {
	return nil
}

func (curM StateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case errorMsg:
		curM.err = msg
		return curM, tea.Quit

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		curM.locationList.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if curM.locationList.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, curM.keys.toggleFilter):
			curM.isFilterOpen = !curM.isFilterOpen
			curM.locationList.SetFilteringEnabled(curM.isFilterOpen)

		case key.Matches(msg, curM.keys.choose):
			i, ok := curM.locationList.SelectedItem().(item)
			if !ok {
				break
			}
			curM.curitem = &i
			curM.isFilterOpen = false
			return curM, nil
		case key.Matches(msg, curM.keys.quitApp):
			return curM, tea.Quit
		}
	}

	if curM.isFilterOpen {
		curM.locationList, cmd = curM.locationList.Update(msg)
	}

	return curM, cmd
}

func (curM StateModel) View() string {

	if curM.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", curM.err)
	}

	if curM.isFilterOpen {
		return appStyle.Render(curM.locationList.View())
	}

	helpView := curM.help.View(curM.keys)
	height := 8 - strings.Count(helpView, "\n")
	if curM.curitem == nil {
		return "Filtering is not rendered" + strings.Repeat("\n", height) + helpView
	}
	return fmt.Sprint(*curM.curitem) + strings.Repeat("\n", height) + helpView
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

	newModel := StateModel{
		locationList: list.New(newItems, list.NewDefaultDelegate(), 0, 0),
		isFilterOpen: false,
		keys:         newItemsKeyMap(),
		help:         help.New(),
	}
	newModel.locationList.Title = "Find Cities"
	newModel.locationList.SetShowFilter(true)
	newModel.locationList.SetFilteringEnabled(false)

	return newModel
}
