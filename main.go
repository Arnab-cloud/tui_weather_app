package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github/Arnab-cloud/tui_weather_app/internal/database"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type location struct {
	Name  string      `json:"name"`
	Coord Coordinates `json:"coord"`
	Id    int         `json:"id"`
}

type apiConfig struct {
	DB *database.Queries
}

var savedCities = make(map[string]City)

func GetCityInfoFromDB(db *database.Queries, loc location) ([]City, error) {
	if loc.Id != 0 {
		data, err := db.FindCityWithID(context.Background(), int64(loc.Id))
		if err != nil {
			return nil, err
		}
		return []City{City{Name: data.Name, Country: data.Country, Lat: data.Lat, Lon: data.Lon, Id: int(data.ID)}}, nil
	} else if loc.Name != "" {
		var cities []City
		query := loc.Name + "%"
		data, err := db.FuzzYFindCity(context.Background(), query)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			return nil, errors.New("no matching city found")
		}

		for _, city := range data {
			cities = append(cities, City{Name: city.Name, Country: city.Country, Lat: city.Lat, Lon: city.Lon, Id: int(city.ID)})
		}

		return cities, nil
	}
	return nil, errors.New("no matching city found")
}

func GetCityInfoFromAPI(ctx context.Context, client *http.Client, loc location, limit int) ([]City, error) {
	baseURL := os.Getenv("GEOCODING_API")
	if baseURL == "" {
		return nil, errors.New("WEATHER_API not set")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, errors.New("API_KEY not set")
	}

	var url string

	if loc.Coord.Lat != 0 && loc.Coord.Lon != 0 {
		url = fmt.Sprintf(
			"%s/reverse?lat=%f&lon=%f&limit=%d&appid=%s",
			baseURL,
			loc.Coord.Lat,
			loc.Coord.Lon,
			limit,
			apiKey,
		)
	} else if loc.Name != "" {
		url = fmt.Sprintf(
			"%s/direct?q=%s&limit=%d&appid=%s",
			baseURL,
			loc.Name,
			limit,
			apiKey,
		)
	} else {
		return nil, errors.New("no valid location city provided")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	cities, err := FetchAndDecode[[]City](client, req)
	if err != nil {
		return nil, err
	}

	if len(*cities) == 0 {
		return nil, errors.New("no cities found")
	}
	return *cities, nil
}

func readCityInfo(city []byte) ([]City, error) {
	var cities []City

	err := json.Unmarshal(city, &cities)
	if err != nil {
		return nil, err
	}
	return cities, err
}

func readCityFromFile() {
	fileName := "sample_res.json"

	city, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error reading response file: %v", err)
	}

	cities, err := readCityInfo(city)

	if err != nil {
		log.Fatalf("error reading json: %v", err)
	}

	if len(cities) == 0 {
		log.Fatal("No cities found")
	}

	for _, city := range cities {
		if _, exists := savedCities[city.Name]; !exists {
			savedCities[city.Name] = city
		}
	}

	fmt.Print(cities[0])
}

func insertBatch(ctx context.Context, db *sql.DB, batch []database.CreateCityParams) error {
	tnx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tnx.Rollback()

	apiCfg := apiConfig{DB: database.New(db)}
	qtx := apiCfg.DB.WithTx(tnx)

	for _, city := range batch {
		if _, err := qtx.CreateCity(ctx, city); err != nil {
			return err
		}
	}

	return tnx.Commit()
}

func loadcityIntoDB(conn *sql.DB) {
	fileName := "city.list.json"
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("Error opening the json file: %v", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if _, err := decoder.Token(); err != nil {

		log.Fatalf("Error readin the first token: %v", err)
	}

	batchSize := 500
	batch := make([]database.CreateCityParams, 0, batchSize)
	count := 0

	for decoder.More() {
		var city struct {
			Country string `json:"country"`
			location
		}

		if err := decoder.Decode(&city); err != nil {
			log.Fatalf("Error reading the docs from file: %v", err)
		}

		log.Printf("city: %v", count)
		count++

		batch = append(batch, database.CreateCityParams{
			ID:      int64(city.Id),
			Name:    city.Name,
			Country: city.Country,
			Lat:     city.Coord.Lat,
			Lon:     city.Coord.Lon,
		})

		if len(batch) == batchSize {
			if err := insertBatch(context.Background(), conn, batch); err != nil {
				log.Fatalf("Error saving the records in db: %v", err)
			}
			batch = batch[:0]
			fmt.Printf("Added %d values", count)
		}
	}

	if len(batch) > 0 {
		if err := insertBatch(context.Background(), conn, batch); err != nil {
			log.Fatalf("Error saving the records in db: %v", err)
		}
	}
}

func readWeatherResponseFile() error {
	fileName := "sample_res.json"
	var weatherRes WeatherResponse

	city, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Error while reading the response file: %v", err)
	}

	if err := json.Unmarshal(city, &weatherRes); err != nil {
		return fmt.Errorf("Error while unmarshalling the weather response: %v", err)
	}

	fmt.Printf("city:\n%v", weatherRes)

	return nil
}

func GetWeatherInfo(db *database.Queries, client *http.Client, loc location) error {
	if loc.Coord.Lat != 0 && loc.Coord.Lon != 0 {
		return GetWeatherInfoFromCoord(db, client, loc.Coord)
	} else if loc.Name != "" {
		if cities, err := GetCityInfoFromAPI(context.Background(), client, loc, 10); err != nil {
			return err
		} else {
			return GetWeatherInfoFromCoord(db, client, Coordinates{Lat: cities[0].Lat, Lon: cities[0].Lon})
		}
	}
	return errors.New("not enough information to get weather info")
}

func GetWeatherInfoFromCoord(db *database.Queries, client *http.Client, coord Coordinates) error {

	lastUpdated := sql.NullInt64{Int64: time.Now().Add(-5 * time.Minute).Unix(), Valid: true}
	query := database.GetFreshWeatherByCoordsParams{
		Lat:       sql.NullFloat64{Float64: coord.Lat, Valid: true},
		Lon:       sql.NullFloat64{Float64: coord.Lon, Valid: true},
		FetchedAt: lastUpdated,
	}
	if city, err := db.GetFreshWeatherByCoords(context.Background(), query); err == nil {
		fmt.Print(city)
		return nil
	}
	log.Printf("No cached value or old cached value")

	city, err := getWeatherInfoAPI(client, coord, 0, "")
	if err != nil {
		return fmt.Errorf("Failed to fetch weather info: %s", err)
	}

	err = db.InsertWeather(context.Background(), city.ToDBWeather())

	fmt.Print(city)
	return err

}

func getWeatherInfoAPI(client *http.Client, coord Coordinates, cityCode int, cityName string) (*WeatherResponse, error) {
	weatherURL := os.Getenv("WEATHER_API")
	if weatherURL == "" {
		return nil, errors.New("no weather url found in the environtment")
	}
	apiKey := os.Getenv("API_KEy")
	if apiKey == "" {
		return nil, errors.New("no api key found in the environtment")
	}

	requestURI := weatherURL
	if (Coordinates{}) != coord {
		requestURI += fmt.Sprintf("?lat=%f&lon=%f", coord.Lat, coord.Lon)
	} else if cityCode != 0 {
		requestURI += fmt.Sprintf("?id=%d", cityCode)
	} else if cityName != "" {
		requestURI += fmt.Sprintf("?q=%s", cityName)
	} else {
		return nil, errors.New("either cityCode or cityName must be non zero")
	}

	requestURI += fmt.Sprintf("&appid=%s", apiKey)

	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		return nil, err
	}

	return FetchAndDecode[WeatherResponse](client, req)
}

func GetCityInfo(db *database.Queries, client *http.Client, loc location) ([]City, error) {
	if cities, err := GetCityInfoFromDB(db, loc); err == nil {
		return cities, nil
	}

	log.Print("city info not found in db")
	log.Print("calling the geocoding api")

	return GetCityInfoFromAPI(context.Background(), client, loc, 10)
}

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
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Db url not found")
	}
	client := &http.Client{}
	defer client.CloseIdleConnections()

	conn, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Error at db connection: %v", err)
	}
	defer conn.Close()

	queries := database.New(conn)

	apiCfg := apiConfig{
		DB: queries,
	}

	cities, err := GetCityInfo(apiCfg.DB, client, location{Name: cityName})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(cities)
}
