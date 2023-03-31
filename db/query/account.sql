-- name: CreateAccount :one
INSERT into accounts(owner,
                     balance,
                     currency)
values ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
select *
from accounts
where id = $1
limit 1;

-- name: GetAccountForUpdate :one
select *
from accounts
where id = $1
limit 1 FOR no key UPDATE;

-- name: ListAccounts :many
select *
from accounts
order by id
limit $1 offset $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
where id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
where id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE
from accounts
where id = $1;