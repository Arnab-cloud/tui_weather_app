package weather

import (
	"database/sql"
	"github/Arnab-cloud/tui_weather_app/internal/database"
	"time"
)

func (res *WeatherResponse) ToDBWeather() database.InsertWeatherParams {
	now := time.Now().Unix()
	var (
		weatherMain sql.NullString
		weatherDesc sql.NullString
		weatherIcon sql.NullString
		weatherID   sql.NullInt64
	)
	// Weather array is usually non-empty, but be safe
	if len(res.Weather) > 0 {
		weatherMain = sql.NullString{String: res.Weather[0].Type, Valid: true}
		weatherDesc = sql.NullString{String: res.Weather[0].Desc, Valid: true}
		weatherIcon = sql.NullString{String: res.Weather[0].Icon, Valid: true}
		weatherID = sql.NullInt64{Int64: int64(res.Weather[0].Id), Valid: true}
	}
	// Rain is optional in OpenWeatherMap
	var rain1h sql.NullFloat64
	if res.Rain > 0 {
		rain1h = sql.NullFloat64{Float64: float64(res.Rain), Valid: true}
	}

	// Country is optional
	var country sql.NullString
	if res.Sys.Country != "" {
		country = sql.NullString{String: res.Sys.Country, Valid: true}
	}

	return database.InsertWeatherParams{
		CityID:      sql.NullInt64{Int64: int64(res.ID), Valid: true},
		CityName:    sql.NullString{String: res.Name, Valid: res.Name != ""},
		Country:     country,
		Lat:         sql.NullFloat64{Float64: res.Coord.Lat, Valid: true},
		Lon:         sql.NullFloat64{Float64: res.Coord.Lon, Valid: true},
		WeatherMain: weatherMain,
		WeatherDesc: weatherDesc,
		WeatherIcon: weatherIcon,
		WeatherID:   weatherID,
		Temp:        sql.NullFloat64{Float64: float64(res.Main.Temp), Valid: true},
		FeelsLike:   sql.NullFloat64{Float64: float64(res.Main.FeelsLike), Valid: true},
		TempMin:     sql.NullFloat64{Float64: float64(res.Main.TempMin), Valid: true},
		TempMax:     sql.NullFloat64{Float64: float64(res.Main.TempMax), Valid: true},
		Humidity:    sql.NullInt64{Int64: int64(res.Main.Humidity), Valid: true},
		Pressure:    sql.NullInt64{Int64: int64(res.Main.Pressure), Valid: true},
		SeaLevel:    sql.NullInt64{Int64: int64(res.Main.SeaLevel), Valid: res.Main.SeaLevel > 0},
		GroundLevel: sql.NullInt64{Int64: int64(res.Main.GroundLevel), Valid: res.Main.GroundLevel > 0},
		WindSpeed:   sql.NullFloat64{Float64: float64(res.Wind.Speed), Valid: true},
		WindDeg:     sql.NullInt64{Int64: int64(res.Wind.Deg), Valid: true},
		WindGust:    sql.NullFloat64{Float64: float64(res.Wind.Gust), Valid: res.Wind.Gust > 0},
		Rain1h:      rain1h,
		Cloudiness:  sql.NullInt64{Int64: int64(res.Clouds), Valid: true},
		Visibility:  sql.NullInt64{Int64: int64(res.Vis), Valid: true},
		Sunrise:     sql.NullInt64{Int64: res.Sys.Sunrise, Valid: res.Sys.Sunrise > 0},
		Sunset:      sql.NullInt64{Int64: res.Sys.Sunset, Valid: res.Sys.Sunset > 0},
		WeatherTime: sql.NullInt64{Int64: res.DT, Valid: true},
		FetchedAt:   sql.NullInt64{Int64: now, Valid: true},
		Timezone:    sql.NullInt64{Int64: int64(res.Timezone), Valid: true},
	}
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

		Sys: WeatherSys{
			Country: nullString(w.Country),
			Sunrise: nullInt64(w.Sunrise),
			Sunset:  nullInt64(w.Sunset),
			Id:      int(w.ID),
		},

		Weather: []BasicWeather{
			{
				Type: nullString(w.WeatherMain),
				Desc: nullString(w.WeatherDesc),
				Icon: nullString(w.WeatherIcon),
				Id:   nullInt(w.WeatherID),
			},
		},

		Main: MainWeather{
			Temp:        nullFloat64(w.Temp),
			FeelsLike:   nullFloat64(w.FeelsLike),
			TempMin:     nullFloat64(w.TempMin),
			TempMax:     nullFloat64(w.TempMax),
			Pressure:    nullInt(w.Pressure),
			Humidity:    nullInt(w.Humidity),
			SeaLevel:    nullInt(w.SeaLevel),
			GroundLevel: nullInt(w.GroundLevel),
		},

		Wind: Wind{
			Speed: nullFloat64(w.WindSpeed),
			Gust:  nullFloat64(w.WindGust),
			Deg:   nullInt(w.WindDeg),
		},

		Rain: nullFloat64(w.Rain1h),
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
