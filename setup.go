package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "embed"
)

var (
	GEOCODING_API string = ""
	WEATHER_API   string = ""
)

//go:embed small.db.gz
var embeddedDBgz []byte

func createUserAppDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(configDir, "tui_weather_app")
	return appDir, nil
}

// ensureDatabase checks if the DB exists in the user's config folder.
// If not, it extracts the gzipped embedded DB to that location.
func ensureDatabase(appDir string) (string, error) {

	dbPath := filepath.Join(appDir, "weather.db")

	// If file exists, return path immediately
	if _, err := os.Stat(dbPath); err == nil {
		return dbPath, nil
	}

	// Create directory if missing
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("could not create app directory: %w", err)
	}

	// Decompress and Write
	log.Println("Extracting embedded database for first-time use...")

	reader, err := gzip.NewReader(bytes.NewReader(embeddedDBgz))
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer reader.Close()

	outFile, err := os.Create(dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to create local db file: %w", err)
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, reader); err != nil {
		return "", fmt.Errorf("failed to extract database content: %w", err)
	}

	log.Println("Database setup complete.")
	return dbPath, nil
}

func GetEnvVariables() (ApiKey, WeatherApi, GeocodingApi string) {
	ApiKey = os.Getenv("API_KEY")
	if ApiKey == "" {
		log.Fatal("API_KEY is required. Set it as an environment variable.")
	}

	if WEATHER_API == "" {
		WeatherApi = os.Getenv("WEATHER_API")
	} else {
		WeatherApi = WEATHER_API
	}

	if GEOCODING_API == "" {
		GeocodingApi = os.Getenv("GEOCODING_API")
	} else {
		GeocodingApi = GEOCODING_API
	}

	if WeatherApi == "" || GeocodingApi == "" {
		log.Fatal("WeatherApi or Geocoding Api is missing")
	}

	return
}
