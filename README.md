# TUI Weather App

A lightweight Terminal User Interface (TUI) weather application built with Go. It features a local SQLite cache to minimize API calls and support fuzzy city searching.

## ✨ Features

* **Cache-Aside Architecture:** Checks local SQLite database before hitting the API.
* **Fuzzy Search:** Look up cities with partial names.
* **Efficient Data Ingestion:** Streams large city datasets into SQLite using JSON decoding and batch transactions.

## 🛠️ Prerequisites

* [Go](https://go.dev/doc/install) (1.21+)
* [OpenWeatherMap API Key](https://openweathermap.org/api)

## ⚙️ Setup

1. **Clone the repository:**
```bash
git clone https://github.com/Arnab-cloud/tui_weather_app.git
cd tui_weather_app

```


2. **Configure Environment Variables:**
Create a `.env` file in the root directory:
```env
API_KEY=your_openweather_api_key_here
WEATHER_API=https://api.openweathermap.org/data/2.5/weather
GEOCODING_API=http://api.openweathermap.org/geo/1.0
DB_URL=./weather.db

```


3. **Initialize Database:**
Ensure you have your `city.list.json` in the root folder, then run the seeder (if implemented in your main) to populate the local database.

## 🚀 Usage

Run the application by providing a city name as an argument:

```bash
go run . Kolkata

```

If no city is provided, it defaults to the city configured in `main.go`.

## 📂 Project Structure

* `/internal/database`: SQLC generated code and DB seeding logic.
* `client.go`: Handles raw HTTP communication with OpenWeather APIs.
* `service.go`: Business logic (Orchestrates DB vs API calls).
* `types.go`: JSON and Database data models.
