
-- +goose Up
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    tag text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE UNIQUE INDEX unique_tag ON categories (LOWER(tag));


-- +goose Down
DROP TABLE categories;
DROP INDEX unique_tag;
