
-- +goose Up
CREATE TABLE invitations (
    id BIGSERIAL PRIMARY KEY,
    category_id bigint NOT NULL REFERENCES categories(id),
    private_id text NOT NULL,
    greeting text NOT NULL,
    maximum_guest_count int NOT NULL DEFAULT 1,
    status text,
    notes text,
    mobile_phone_number text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE UNIQUE INDEX unique_private_id ON invitations (private_id);
CREATE UNIQUE INDEX unique_greeting ON invitations (LOWER(greeting));
CREATE UNIQUE INDEX unique_mobile_phone_number ON invitations (LOWER(mobile_phone_number));


-- +goose Down
DROP TABLE invitations;
