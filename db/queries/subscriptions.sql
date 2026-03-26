-- name: CreateSubscription :one
insert into subscriptions (
    service_name, price, user_id, start_date, end_date
) values (
    $1, $2, $3, $4, $5
)
returning *;

-- name: ReadSubscription :one
select *
from subscriptions
where id = $1;

-- name: UpdateSubscription :one
update subscriptions
set service_name = $2,
    price = $3,
    user_id = $4,
    start_date = $5,
    end_date = $6
where id = $1
returning *;

-- name: DeleteSubscription :one
delete from subscriptions
where id = $1
returning *;

-- name: ListSubscriptions :many
select *
from subscriptions
where (sqlc.narg('user_id')::uuid is null or user_id = sqlc.narg('user_id')::uuid)
    and (sqlc.narg('price')::int is null or price = sqlc.narg('price')::int)
    and (sqlc.narg('service_name')::text is null or service_name = sqlc.narg('service_name')::text)
    and (sqlc.narg('start_date')::date is null or start_date = sqlc.narg('start_date')::date)
    and (sqlc.narg('end_date')::date is null or end_date = sqlc.narg('end_date')::date);

-- name: TotalSubscriptionsCost :one
select coalesce(sum(
    (
        (
            extract(year from overlap_end) * 12 + extract(month from overlap_end)
        ) -
        (
            extract(year from overlap_start) * 12 + extract(month from overlap_start)
        ) + 1
    ) * price
), 0)::bigint AS cost
from (
    select
        price,
        greatest(start_date, sqlc.arg('start_date')::date) as overlap_start,
        least(coalesce(end_date, sqlc.arg('end_date')::date), sqlc.arg('end_date')::date) as overlap_end
    from subscriptions
    where
        start_date <= sqlc.arg('end_date')::date
        and (end_date is null or end_date >= sqlc.arg('start_date')::date)
        and (sqlc.narg('user_id')::uuid is null or user_id = sqlc.narg('user_id')::uuid)
        and (sqlc.narg('service_name')::text is null or service_name = sqlc.narg('service_name')::text)
) t;
