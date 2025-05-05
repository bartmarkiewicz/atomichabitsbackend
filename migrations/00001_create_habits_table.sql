-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS habits (
    id UUID PRIMARY KEY,
    description TEXT,
    icon_base64 TEXT,
    colour_hex TEXT,
    mode_type TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS habits;
-- +goose StatementEnd
