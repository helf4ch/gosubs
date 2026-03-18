-- +goose Up
-- +goose StatementBegin
create table subscriptions (
    id uuid primary key default uuidv7(),
    service_name text not null,
    price int not null,
    user_id uuid not null,
    start_date date not null,
    end_date date,
    created_at timestamp not null default now(),
    updated_at timestamp,
    constraint unique_sub_info unique (service_name, user_id, start_date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists subscriptions;
-- +goose StatementEnd
