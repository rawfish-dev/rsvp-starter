
-- +goose Up
DROP INDEX unique_mobile_phone_number;

-- +goose Down
CREATE UNIQUE INDEX unique_mobile_phone_number ON invitations (LOWER(mobile_phone_number));

