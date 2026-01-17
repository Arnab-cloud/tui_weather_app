-- +goose Up
CREATE INDEX city_name_idx ON cities (name);

-- +goose Down
DROP INDEX city_name_idx FROM cities;
