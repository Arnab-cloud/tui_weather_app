-- +goose Up
CREATE INDEX idx_weather_cid ON weather_cache(city_id);
CREATE INDEX idx_weather_city ON weather_cache(city_name, country);
CREATE INDEX idx_weather_coords ON weather_cache(lat, lon);
CREATE INDEX idx_weather_fetched ON weather_cache(fetched_at);

-- +goose Down
DROP INDEX idx_weather_fetched FROM weather_cache;
DROP INDEX idx_weather_coords FROM weather_cache;
DROP INDEX idx_weather_city FROM weather_cache;
DROP INDEX idx_weather_cid FROM weather_cache;
