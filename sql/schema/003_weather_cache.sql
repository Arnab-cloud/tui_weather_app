-- +goose Up
CREATE TABLE weather_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    -- Location
    city_id INTEGER,
    city_name TEXT,
    country TEXT,
    lat REAL,
    lon REAL,

    -- Weather summary
    weather_main TEXT,
    weather_desc TEXT,
    weather_icon TEXT,

    -- Temperature
    temp REAL,
    feels_like REAL,
    temp_min REAL,
    temp_max REAL,
    humidity INTEGER,
    pressure INTEGER,

    -- Wind
    wind_speed REAL,
    wind_deg INTEGER,
    wind_gust REAL,

    -- Rain (nullable)
    rain_1h REAL,

    -- Clouds
    cloudiness INTEGER,

    -- Visibility
    visibility INTEGER,

    -- Time
    weather_time INTEGER,   -- `dt` from API
    fetched_at INTEGER,    -- when YOU fetched it (unix time)
    timezone INTEGER
);

-- +goose Down
DROP TABLE weather_cache;
