-- +goose Up
-- +goose StatementBegin
alter table ozon.orders add column weight_in_kg float4 not null default -1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table ozon.orders drop column weight_in_kg;
-- +goose StatementEnd
