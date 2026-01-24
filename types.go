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
	Temp        float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	TempMin     float64 `json:"temp_min"`
	TempMax     float64 `json:"temp_max"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
	SeaLevel    int     `json:"sea_level"`
	GroundLevel int     `json:"grnd_level"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Gust  float64 `json:"gust"`
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
	Sys      WeatherSys     `json:"sys"`
	Wind     Wind           `json:"wind"`
	Coord    Coordinates    `json:"coord"`
	Rain     float64        `json:"rain.1h"`
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

type Location struct {
	Name  string      `json:"name"`
	Coord Coordinates `json:"coord"`
	Id    int         `json:"id"`
}

func WeatherCacheToResponse(w database.WeatherCache) WeatherResponse {
	return WeatherResponse{
		ID:   int(w.ID),
		Name: nullString(w.CityName),
		DT:   nullInt64(w.WeatherTime),
		COD:  200, // cached result is assumed OK
		Base: "stations",

		Timezone: nullInt(w.Timezone),
		Vis:      nullInt(w.Visibility),
		Clouds:   nullInt(w.Cloudiness),

		Coord: Coordinates{
			Lat: nullFloat64(w.Lat),
			Lon: nullFloat64(w.Lon),
		},

		Weather: []BasicWeather{
			{
				Type: nullString(w.WeatherMain),
				Desc: nullString(w.WeatherDesc),
				Icon: nullString(w.WeatherIcon),
				Id:   0, // not present in cache
			},
		},

		Main: MainWeather{
			Temp:        nullFloat64(w.Temp),
			FeelsLike:   nullFloat64(w.FeelsLike),
			TempMin:     nullFloat64(w.TempMin),
			TempMax:     nullFloat64(w.TempMax),
			Pressure:    nullInt(w.Pressure),
			Humidity:    nullInt(w.Humidity),
			SeaLevel:    0, // not stored
			GroundLevel: 0, // not stored
		},

		Wind: Wind{
			Speed: nullFloat64(w.WindSpeed),
			Gust:  nullFloat64(w.WindGust),
			Deg:   nullInt(w.WindDeg),
		},

		Rain: nullFloat64(w.Rain1h),
	}
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

func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullInt64(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

func nullInt(ni sql.NullInt64) int {
	if ni.Valid {
		return int(ni.Int64)
	}
	return 0
}

func nullFloat64(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}
