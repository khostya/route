-- +goose Up
-- +goose StatementBegin
create table if not exists ozon.orders
(
    id                text primary key         not null,
    recipient_id      text                     not null,
    status            ozon.status              not null,
    status_updated_at timestamp with time zone not null,
    hash              text                     not null,
    created_at        timestamp with time zone not null,
    expiration_date   timestamp with time zone not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ozon.orders;
-- +goose StatementEnd
