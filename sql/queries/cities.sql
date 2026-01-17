-- name: CreateCity :one
INSERT INTO cities (id, name, country, lat, lon)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteCity :exec
DELETE FROM cities
WHERE id = ?;
