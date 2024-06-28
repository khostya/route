-- +goose Up
-- +goose StatementBegin
create type ozon.wrapper as enum ('box', 'package', 'stretch');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists ozon.wrapper;
-- +goose StatementEnd
