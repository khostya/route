-- +goose Up
-- +goose StatementBegin
alter table ozon.orders add column price_in_rub numeric not null default 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table ozon.orders drop column price_in_rub;
-- +goose StatementEnd
