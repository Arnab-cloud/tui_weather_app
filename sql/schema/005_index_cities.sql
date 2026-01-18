-- +goose Up
DROP INDEX city_name_idx;
CREATE INDEX city_name_idx ON cities (name COLLATE NOCASE);

-- +goose Down
DROP INDEX city_name_idx;
CREATE INDEX city_name_idx ON cities (name);
