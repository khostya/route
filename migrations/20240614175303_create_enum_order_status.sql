-- +goose Up
-- +goose StatementBegin
create type ozon.status as enum ('delivered', 'issued', 'refunded');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists ozon.status;
-- +goose StatementEnd
