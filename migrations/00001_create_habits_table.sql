-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS habits (
    id UUID PRIMARY KEY,
    description TEXT,
    icon_base_64 TEXT,
    colourHex TEXT,
    mode_type TEXT
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS habits;
-- +goose StatementEnd
