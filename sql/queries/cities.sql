-- name: CreateCity :one
INSERT INTO cities (id, name, country, lat, lon)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteCity :exec
DELETE FROM cities
WHERE id = ?;

-- name: FindCity :many
SELECT *
FROM cities
WHERE LOWER(name)=LOWER(?);


-- name: FindCityWithID :one
SELECT *
FROM cities
WHERE id = ?;

-- name: FuzzYFindCity :many
SELECT *
FROM cities
WHERE name like ?;
