-- +goose Up
-- +goose StatementBegin
update users set role=0 where id=1
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
update users set role=2 where id=1
-- +goose StatementEnd
