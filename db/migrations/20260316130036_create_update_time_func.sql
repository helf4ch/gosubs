-- +goose Up
-- +goose StatementBegin
create function update_time() returns trigger as $$
begin
    new.updated_at := now();
    return new;
end;
$$ language plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop function update_time();
-- +goose StatementEnd
