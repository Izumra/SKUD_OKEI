-- +goose Up
-- +goose StatementBegin
update users set role=1 where id=2
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update users set role=1 where id=2
-- +goose StatementEnd
