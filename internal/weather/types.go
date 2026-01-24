package weather

import (
	"fmt"
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

type Location struct {
	Name  string      `json:"name"`
	Coord Coordinates `json:"coord"`
	Id    int         `json:"id"`
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

func (city City) Title() string { return city.Name }
func (city City) Description() string {
	return fmt.Sprintf("%s, Lat: %f, Lon: %f", city.Country, city.Lat, city.Lon)
}
func (city City) FilterValue() string { return city.Name }
