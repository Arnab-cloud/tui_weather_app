package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github/Arnab-cloud/tui_weather_app/internal/database"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type City struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Id      int     `json:"id"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type rawCity struct {
	Name    string      `json:"name"`
	Country string      `json:"country"`
	Coord   Coordinates `json:"coord"`
	Id      int         `json:"id"`
}

type SavedCities struct {
	Cities []City `json:"cities"`
}

type apiConfig struct {
	DB *database.Queries
}

var savedCities = make(map[string]City)

func getCityInfo(cityName string, limit int) (City, error) {
	err := godotenv.Load()
	if err != nil {
		return City{}, fmt.Errorf("Error loading env variables: %s", err)
	}

	Weather_Api := os.Getenv("WEATHER_API")
	if Weather_Api == "" {
		return City{}, fmt.Errorf("Error loading Weather Api stirng: %s", err)
	}

	Api_Key := os.Getenv("API_KEY")
	if Api_Key == "" {
		return City{}, fmt.Errorf("Error loading Api_key Api stirng: %s", err)
	}

	requestUrl := fmt.Sprintf("%s/direct?q=%s&limit=%d&appid=%s", Weather_Api, cityName, limit, Api_Key)
	resp, err := http.Get(requestUrl)
	if err != nil {
		return City{}, fmt.Errorf("Error in getting city info: %s", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return City{}, fmt.Errorf("Error reading city info api response: %s", err)
	}

	resp.Body.Close()

	if resp.StatusCode > 299 {
		return City{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}

	cities, err := readCityInfo(body)
	if err != nil {
		return City{}, err
	}

	if len(cities) == 0 {
		return City{}, errors.New("No cities found")
	}

	return cities[0], nil
}

func readCityInfo(data []byte) ([]City, error) {
	var cities []City

	err := json.Unmarshal(data, &cities)
	if err != nil {
		return nil, err
	}
	return cities, err
}

func readCityFromFile() {
	fileName := "sample_res.json"

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error reading response file: %v", err)
	}

	cities, err := readCityInfo(data)
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

func loadDataIntoDB(conn *sql.DB) {
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
		var city rawCity

		if err := decoder.Decode(&city); err != nil {
			log.Fatalf("Error reading the docs from file: %v", err)
		}

		log.Printf("data: %v", count)
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

func main() {
	// city, err := getCityInfo("Kolkata", 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// readCityFromFile()
	// data, err := getCityInfo("kolkata", 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// readCityFromFile()
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Db url not found")
	}

	conn, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		log.Fatalf("Error at db connection: %v", err)
	}
	defer conn.Close()

	// queries := database.New(conn)

	// apiCfg := apiConfig{
	// 	DB: queries,
	// }

	// // data := database.CreateCityParams{ID: 1000000, Name: "first city", Country: "first country", Lat: 12.234, Lon: 233.334}

	// // dbCitym, err := apiCfg.DB.CreateCity(context.Background(), data)
	// err = apiCfg.DB.DeleteCity(context.Background(), 1000000)
	// if err != nil {
	// 	log.Fatalf("Error executing the query: %v", err)
	// }

	// fmt.Print(dbRes)
	loadDataIntoDB(conn)
}
