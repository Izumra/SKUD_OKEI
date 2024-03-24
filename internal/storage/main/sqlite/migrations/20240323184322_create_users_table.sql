-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    pass VARCHAR(25) NOT NULL,
    role INTEGER NOT NULL DEFAULT 2  
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
