-- name: GetLatestWeatherByCity :one
SELECT *
FROM weather_cache
WHERE city_name = ?
ORDER BY fetched_at DESC
LIMIT 1;

-- name: GetFreshWeatherByCity :one
SELECT *
FROM weather_cache
WHERE city_name = ?
  AND fetched_at >= ?
ORDER BY fetched_at DESC
LIMIT 1;


-- name: GetLatestWeatherByCoords :one
SELECT *
FROM weather_cache
WHERE lat = ?
  AND lon = ?
ORDER BY fetched_at DESC
LIMIT 1;


-- name: GetFreshWeatherByCoords :one
SELECT *
FROM weather_cache
WHERE lat = ?
  AND lon = ?
  AND fetched_at >= ?
ORDER BY fetched_at DESC
LIMIT 1;


-- name: GetLatestWeatherByCityID :one
SELECT *
FROM weather_cache
WHERE city_id = ?
ORDER BY fetched_at DESC
LIMIT 1;


-- name: InsertWeather :exec
INSERT INTO weather_cache (
    city_id,
    city_name,
    country,
    lat,
    lon,
    weather_main,
    weather_desc,
    weather_icon,
    temp,
    feels_like,
    temp_min,
    temp_max,
    humidity,
    pressure,
    wind_speed,
    wind_deg,
    wind_gust,
    rain_1h,
    cloudiness,
    visibility,
    weather_time,
    fetched_at,
    timezone
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);



-- name: DeleteOldWeather :exec
DELETE FROM weather_cache
WHERE fetched_at < ?;


-- name: DeleteDuplicateWeather :exec
DELETE FROM weather_cache
WHERE id NOT IN (
    SELECT MAX(id)
    FROM weather_cache
    GROUP BY city_id
);



-- name: GetWeatherHistoryByCity :many
SELECT *
FROM weather_cache
WHERE city_name = ?
ORDER BY fetched_at DESC
LIMIT ?;
