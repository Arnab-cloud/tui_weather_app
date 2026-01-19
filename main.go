package main

import (
	"context"
	"database/sql"
	"fmt"
	"github/Arnab-cloud/tui_weather_app/internal/database"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	cityName := "Kolkata"
	commandArgs := os.Args[1:]
	if len(commandArgs) == 0 {
		log.Printf("No city name provided, using default: %s", cityName)
	} else {
		cityName = strings.Join(commandArgs, " ")
		log.Printf("Loooking up weather for city: %s", cityName)
	}

	godotenv.Load()

	conn, err := sql.Open("sqlite3", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error at db connection: %v", err)
	}
	defer conn.Close()

	client := NewWeatherClient(
		os.Getenv("API_KEY"),
		os.Getenv("WEATHER_API"),
		os.Getenv("GEOCODING_API"),
	)

	service := &WeatherService{
		DB:     database.New(conn),
		Client: client,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	weather, err := service.GetWeather(ctx, Location{Name: cityName})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(*weather)
}
