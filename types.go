package main

import (
	"database/sql"
	"github/Arnab-cloud/tui_weather_app/internal/database"
	"time"
)

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type BasicWeather struct {
	Type string `json:"main"`
	Desc string `json:"description"`
	Icon string `json:"icon"`
	Id   int    `json:"id"`
}

type MainWeather struct {
	Temp        float32 `json:"temp"`
	FeelsLike   float32 `json:"feels_like"`
	TempMin     float32 `json:"temp_min"`
	TempMax     float32 `json:"temp_max"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
	SeaLevel    int     `json:"sea_level"`
	GroundLevel int     `json:"grnd_level"`
}

type Wind struct {
	Speed float32 `json:"speed"`
	Gust  float32 `json:"gust"`
	Deg   int     `json:"deg"`
}

type WeatherSys struct {
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
	Type    int    `json:"type"`
	Id      int    `json:"id"`
}

type WeatherResponse struct {
	Weather  []BasicWeather `json:"weather"`
	Main     MainWeather    `json:"main"`
	Wind     Wind           `json:"wind"`
	Coord    Coordinates    `json:"coord"`
	Rain     float32        `json:"rain.1h"`
	Base     string         `json:"base"`
	Name     string         `json:"name"`
	DT       int64          `json:"dt"`
	COD      int            `json:"cod"`
	ID       int            `json:"id"`
	Clouds   int            `json:"clouds.all"`
	Timezone int            `json:"timezone"`
	Vis      int            `json:"visibility"`
}

type City struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Id      int     `json:"id"`
}

func (res *WeatherResponse) ToDBWeather() database.InsertWeatherParams {
	now := time.Now().Unix()

	var (
		weatherMain sql.NullString
		weatherDesc sql.NullString
		weatherIcon sql.NullString
	)

	// Weather array is usually non-empty, but be safe
	if len(res.Weather) > 0 {
		weatherMain = sql.NullString{String: res.Weather[0].Type, Valid: true}
		weatherDesc = sql.NullString{String: res.Weather[0].Desc, Valid: true}
		weatherIcon = sql.NullString{String: res.Weather[0].Icon, Valid: true}
	}

	// Rain is optional in OpenWeatherMap
	var rain1h sql.NullFloat64
	if res.Rain > 0 {
		rain1h = sql.NullFloat64{Float64: float64(res.Rain), Valid: true}
	}

	return database.InsertWeatherParams{
		CityID:   sql.NullInt64{Int64: int64(res.ID), Valid: true},
		CityName: sql.NullString{String: res.Name, Valid: res.Name != ""},

		Lat: sql.NullFloat64{Float64: res.Coord.Lat, Valid: true},
		Lon: sql.NullFloat64{Float64: res.Coord.Lon, Valid: true},

		WeatherMain: weatherMain,
		WeatherDesc: weatherDesc,
		WeatherIcon: weatherIcon,

		Temp:      sql.NullFloat64{Float64: float64(res.Main.Temp), Valid: true},
		FeelsLike: sql.NullFloat64{Float64: float64(res.Main.FeelsLike), Valid: true},
		TempMin:   sql.NullFloat64{Float64: float64(res.Main.TempMin), Valid: true},
		TempMax:   sql.NullFloat64{Float64: float64(res.Main.TempMax), Valid: true},

		Humidity: sql.NullInt64{Int64: int64(res.Main.Humidity), Valid: true},
		Pressure: sql.NullInt64{Int64: int64(res.Main.Pressure), Valid: true},

		WindSpeed: sql.NullFloat64{Float64: float64(res.Wind.Speed), Valid: true},
		WindDeg:   sql.NullInt64{Int64: int64(res.Wind.Deg), Valid: true},
		WindGust:  sql.NullFloat64{Float64: float64(res.Wind.Gust), Valid: res.Wind.Gust > 0},

		Rain1h: rain1h,

		Cloudiness: sql.NullInt64{Int64: int64(res.Clouds), Valid: true},
		Visibility: sql.NullInt64{Int64: int64(res.Vis), Valid: true},

		WeatherTime: sql.NullInt64{Int64: res.DT, Valid: true},
		FetchedAt:   sql.NullInt64{Int64: now, Valid: true},
		Timezone:    sql.NullInt64{Int64: int64(res.Timezone), Valid: true},
	}
}
