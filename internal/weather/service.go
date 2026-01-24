package weather

import (
	"context"
	"database/sql"
	"fmt"
	"github/Arnab-cloud/tui_weather_app/internal/database"
	"log"
	"time"
)

type WeatherService struct {
	DB     *database.Queries
	Client *WeatherClient
}

const (
	EPSILON       = 0.01
	CacheDuration = 10 * time.Minute
)

func NewWeatherService(conn *sql.DB, client *WeatherClient) *WeatherService {
	return &WeatherService{
		DB:     database.New(conn),
		Client: client,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, loc Location) (*WeatherResponse, error) {
	var w *WeatherResponse
	if loc.Coord.Lat == 0 && loc.Coord.Lon == 0 {
		cities, err := s.ResolveCity(ctx, loc.Name)
		if err != nil {
			return nil, err
		}
		loc.Coord = Coordinates{Lat: cities[0].Lat, Lon: cities[0].Lon}
	}

	cacheParams := database.GetFreshWeatherByCoordsParams{
		Lat:       sql.NullFloat64{Float64: loc.Coord.Lat - EPSILON, Valid: true},
		Lat_2:     sql.NullFloat64{Float64: loc.Coord.Lat + EPSILON, Valid: true},
		Lon:       sql.NullFloat64{Float64: loc.Coord.Lon - EPSILON, Valid: true},
		Lon_2:     sql.NullFloat64{Float64: loc.Coord.Lon + EPSILON, Valid: true},
		FetchedAt: sql.NullInt64{Int64: time.Now().Add(-10 * time.Minute).Unix(), Valid: true},
	}

	if cached, err := s.DB.GetFreshWeatherByCoords(ctx, cacheParams); err == nil {
		cachedWeatherRes := WeatherCacheToResponse(cached)
		return &cachedWeatherRes, nil
	}

	w, err := s.Client.FetchWeather(ctx, loc.Coord.Lat, loc.Coord.Lon)
	if err != nil {
		return nil, err
	}

	func() {
		if err := s.DB.InsertWeather(ctx, w.ToDBWeather()); err != nil {
			log.Printf("Failed to cache the response: %s", err)
		}
	}()

	return w, nil
}

func (s *WeatherService) ResolveCity(ctx context.Context, name string) ([]City, error) {
	var cities []City
	query := name + "%"
	dbCities, err := s.DB.FuzzYFindCity(ctx, query)

	log.Printf("db queried")

	if err != nil {
		log.Printf("fuzzyfind Error: %v", err)
	}

	if err == nil && len(dbCities) > 0 {
		for _, dbCity := range dbCities {
			cities = append(cities, City{
				Id:      int(dbCity.ID),
				Name:    dbCity.Name,
				Country: dbCity.Country,
				Lat:     dbCity.Lat,
				Lon:     dbCity.Lon,
			})
		}
		return cities, nil
	}

	cities, err = s.Client.FetchGeocoding(ctx, name, 1)
	log.Printf("api queried")
	if err != nil || len(cities) == 0 {
		log.Printf("city '%s' not found locally or via API", name)
		return nil, fmt.Errorf("city '%s' not found locally or via API", name)
	}

	return cities, nil
}
