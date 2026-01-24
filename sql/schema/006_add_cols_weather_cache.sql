-- +goose Up

-- Add weather condition ID
ALTER TABLE weather_cache ADD COLUMN weather_id INTEGER;

-- Add atmospheric pressure readings
ALTER TABLE weather_cache ADD COLUMN sea_level INTEGER;
ALTER TABLE weather_cache ADD COLUMN ground_level INTEGER;

-- Add sun times (essential for TUI display)
ALTER TABLE weather_cache ADD COLUMN sunrise INTEGER;
ALTER TABLE weather_cache ADD COLUMN sunset INTEGER;

-- +goose Down

ALTER TABLE weather_cache DROP COLUMN sunset INTEGER;
ALTER TABLE weather_cache DROP COLUMN sunrise INTEGER;
ALTER TABLE weather_cache DROP COLUMN ground_level INTEGER;
ALTER TABLE weather_cache DROP COLUMN sea_level INTEGER;
ALTER TABLE weather_cache DROP COLUMN weather_id INTEGER;
