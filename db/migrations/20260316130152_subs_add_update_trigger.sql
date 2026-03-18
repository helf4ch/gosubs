-- +goose Up
-- +goose StatementBegin
create trigger subscriptions_update_time before update on subscriptions
for each row execute procedure update_time();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger subscriptions_update_time on subscriptions;
-- +goose StatementEnd
