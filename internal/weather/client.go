package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type WeatherClient struct {
	HTTPClient  *http.Client
	WeatherURL  string
	GeocoderURL string
	APIKey      string
}

func NewWeatherClient(apiKey, weatherURL, geocodeURL string) *WeatherClient {
	return &WeatherClient{
		HTTPClient:  &http.Client{Timeout: 10 * time.Minute},
		WeatherURL:  weatherURL,
		GeocoderURL: geocodeURL,
		APIKey:      apiKey,
	}
}

func FetchAndDecode[T any](client *http.Client, req *http.Request) (*T, error) {
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("Api Error: status=%d, body=%s", response.StatusCode, body)
	}

	var result T
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *WeatherClient) FetchWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error) {
	weatherUrl, err := url.Parse(c.WeatherURL)
	if err != nil {
		return nil, err
	}

	query := weatherUrl.Query()
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("appid", c.APIKey)
	query.Set("units", "metric")

	weatherUrl.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, "GET", weatherUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error forming the weather request: %s", err)
	}

	return FetchAndDecode[WeatherResponse](c.HTTPClient, req)
}

func (c *WeatherClient) FetchGeocoding(ctx context.Context, cityName string, limit int) ([]City, error) {
	geocoderUrl, err := url.Parse(fmt.Sprintf("%s/direct", c.GeocoderURL))
	if err != nil {
		return nil, err
	}

	query := geocoderUrl.Query()
	query.Set("q", cityName)
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("appid", c.APIKey)
	geocoderUrl.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", geocoderUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error forming the geocoding request: %s", err)
	}

	cities, err := FetchAndDecode[[]City](c.HTTPClient, req)
	if err != nil {
		return nil, err
	}

	if len(*cities) == 0 {
		return nil, fmt.Errorf("no cities found with name: %s", cityName)
	}

	return *cities, nil
}

func (c *WeatherClient) FetchReverseGeocoding(ctx context.Context, coord Coordinates, limit int) ([]City, error) {
	geocoderUrl, err := url.Parse(fmt.Sprintf("%s/reverse", c.GeocoderURL))
	if err != nil {
		return nil, err
	}

	query := geocoderUrl.Query()
	query.Set("lat", fmt.Sprintf("%f", coord.Lat))
	query.Set("lon", fmt.Sprintf("%f", coord.Lon))
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("appid", c.APIKey)
	geocoderUrl.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", geocoderUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Error forming the geocoding request: %s", err)
	}

	cities, err := FetchAndDecode[[]City](c.HTTPClient, req)
	if err != nil {
		return nil, err
	}

	if len(*cities) == 0 {
		return nil, fmt.Errorf("no cities found with coordinates: %v", coord)
	}

	return *cities, nil
}
