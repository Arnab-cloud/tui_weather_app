package main

import (
	"database/sql"
	"github/Arnab-cloud/tui_weather_app/internal/ui"
	"github/Arnab-cloud/tui_weather_app/internal/weather"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var Service *weather.WeatherService

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

	client := weather.NewWeatherClient(
		os.Getenv("API_KEY"),
		os.Getenv("WEATHER_API"),
		os.Getenv("GEOCODING_API"),
	)

	service := weather.NewWeatherService(conn, client)

	if _, err := tea.NewProgram(ui.NewModel(service), tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
