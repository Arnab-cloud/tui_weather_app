package ui

import (
	"github/Arnab-cloud/tui_weather_app/internal/weather"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type StateModel struct {
	service *weather.WeatherService

	textInput         textinput.Model
	searchResults     list.Model
	curItem           *weather.City
	curWeather        *weather.WeatherResponse
	isFilterOpen      bool
	isFetchingWeather bool
	keys              *itemsKeyMap
	help              help.Model
	err               error
	debounceId        int
	width             int
	height            int
}

func (curM StateModel) Init() tea.Cmd {
	return nil
}

type citySearchResultMsg struct {
	locs []list.Item
}

type weatherSearchResultMsg struct {
	weather *weather.WeatherResponse
}

type debouncedMsg struct {
	id    int
	query string
}

type errorMsg struct {
	err error
}

func (e errorMsg) Error() string { return e.err.Error() }

func NewModel(service *weather.WeatherService) StateModel {
	newSearchResults := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	ti := textinput.New()
	ti.Placeholder = "Search for a city"
	ti.CharLimit = 50
	ti.Width = 40

	newModel := StateModel{
		service:           service,
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
