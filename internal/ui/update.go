package ui

import (
	"context"
	"github/Arnab-cloud/tui_weather_app/internal/weather"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var debounceDuration = 500 * time.Millisecond

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

		case key.Matches(msg, curM.keys.quit) && !curM.textInput.Focused():
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
			if i, ok := curM.searchResults.SelectedItem().(weather.City); ok {
				curM.curItem = &i
				curM.isFilterOpen = false
				curM.textInput.Blur()
				curM.isFetchingWeather = true
				cmd = curM.performWeatherSearch()
			}
			return curM, cmd
		case key.Matches(msg, curM.keys.back):
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

func (curM *StateModel) debouncedSearch() tea.Cmd {
	query := curM.textInput.Value()
	curM.debounceId++
	return tea.Tick(debounceDuration, func(_ time.Time) tea.Msg {
		return debouncedMsg{id: curM.debounceId, query: query}
	})
}

func (curM StateModel) performLocationSearch() tea.Cmd {
	query := curM.textInput.Value()

	return func() tea.Msg {
		if len(query) < 3 {
			return citySearchResultMsg{locs: nil}
		}
		cities, err := curM.service.ResolveCity(context.Background(), query)
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
		coord := weather.Coordinates{Lat: curM.curItem.Lat, Lon: curM.curItem.Lon}
		res, err := curM.service.GetWeather(
			context.Background(),
			weather.Location{Name: curM.curItem.Name, Coord: coord, Id: curM.curItem.Id},
		)
		log.Print("called get weather")
		if err != nil {
			log.Printf("error fetching the weather: %s", err)
			return weatherSearchResultMsg{weather: nil}
		}
		return weatherSearchResultMsg{weather: res}
	}
}
