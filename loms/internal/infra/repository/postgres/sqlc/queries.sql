-- name: AddOrder :one
insert into orders(user_id, status)
values ($1, $2)
returning order_id;

-- name: GetOrderByID :one
select *
from orders
where order_id = $1;

-- name: UpdateStatusByID :exec
update orders
set status = $2
where order_id = $1;



-- name: AddOrderItem :exec
insert into order_items(sku, order_id, count)
values ($1, $2, $3);

-- name: GetOrderItemsOrderBySKU :many
select sku, order_id, count
from order_items
where order_id = $1
order by sku;



-- name: AddStock :exec
insert into stocks(sku, total_count, reserved)
values ($1, $2, $3)
on conflict (sku)
do update
set total_count = stocks.total_count + $2,
    reserved    = stocks.reserved + $3;

-- name: Reserve :exec
update stocks
set reserved = reserved + $2
where sku = $1;

-- name: RemoveReserve :exec
update stocks
set reserved = reserved - $2
where sku = $1;

-- name: ReduceTotalAndReserve :exec
update stocks
set reserved    = reserved - $2,
    total_count = total_count - $2
where sku = $1;

-- name: GetStockBySKU :one
select *
from stocks
where sku = $1;
