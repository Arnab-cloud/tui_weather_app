package main

import (
	"database/sql"
	"github/Arnab-cloud/tui_weather_app/internal/ui"
	"github/Arnab-cloud/tui_weather_app/internal/weather"
	"log"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	appDir, err := createUserAppDir()
	if err != nil {
		log.Fatalf("Error creating user app dir: %s", err)
	}
	_ = godotenv.Load("release.env", ".env", filepath.Join(appDir, "weather.env"))

	LogFile, err := tea.LogToFile("logs.log", "debug")
	if err != nil {
		log.Printf("Error opening the logfile: %s", err)
		log.Print("Logging disabled")
	}
	defer LogFile.Close()

	dbPath, err := ensureDatabase(appDir)
	if err != nil {
		log.Fatalf("Critical error setting up database: %v", err)
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Error opening the logfile: %s", err)
	}
	defer conn.Close()

	client := weather.NewWeatherClient(GetEnvVariables())

	service := weather.NewWeatherService(conn, client)

	if _, err := tea.NewProgram(ui.NewModel(service), tea.WithAltScreen()).Run(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
