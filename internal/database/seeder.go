package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

type CitySeeder struct {
	db      *sql.DB
	queries *Queries
}

type cityEntry struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Coord   struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"coord"`
}

func NewCitySeeder(db *sql.DB) *CitySeeder {
	return &CitySeeder{
		db:      db,
		queries: New(db),
	}
}

func (s *CitySeeder) LoadCitiesFromFile(ctx context.Context, filePath string, batchSize int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Error opening city file %s", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if _, err := decoder.Token(); err != nil {
		return fmt.Errorf("Error reading JSON start token: %s", err)
	}

	batch := make([]CreateCityParams, 0, batchSize)
	count := 0

	for decoder.More() {
		var entry cityEntry

		if err := decoder.Decode(&entry); err != nil {
			return fmt.Errorf("Decode Error at record %d: %w", count, err)
		}

		batch = append(batch, CreateCityParams{
			ID:      entry.ID,
			Name:    entry.Name,
			Country: entry.Country,
			Lat:     entry.Coord.Lat,
			Lon:     entry.Coord.Lon,
		})

		if len(batch) >= batchSize {
			if err := s.insertBatch(ctx, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
		count++
	}

	if len(batch) > 0 {
		return s.insertBatch(ctx, batch)
	}

	return nil
}

func (s *CitySeeder) insertBatch(ctx context.Context, batch []CreateCityParams) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	for _, city := range batch {
		if _, err := qtx.CreateCity(ctx, city); err != nil {
			return fmt.Errorf("failed to insert city %s: %w", city.Name, err)
		}
	}

	return tx.Commit()
}
