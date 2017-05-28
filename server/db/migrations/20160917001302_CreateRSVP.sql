
-- +goose Up
CREATE TABLE rsvps (
    id BIGSERIAL PRIMARY KEY,
    invitation_private_id text,
    full_name text,
    attending boolean NOT NULL DEFAULT false,
    guest_count int NOT NULL DEFAULT 1,
    special_diet boolean NOT NULL DEFAULT false,
    remarks text,
    mobile_phone_number text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);
CREATE UNIQUE INDEX unique_invitation_private_id ON rsvps (invitation_private_id);

-- +goose Down
DROP TABLE rsvps;
