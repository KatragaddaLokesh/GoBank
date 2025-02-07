-- name: CreateAccount :one
Insert Into account(
owner,
balance,
currency
	) values($1,$2,$3)
returning *;

-- name: GetAccount :one
Select * From account
Where id = $1 Limit 1;

-- name: GetAccountForUpdate :one
Select * From account
Where id = $1 Limit 1 
For no Key Update;


-- name: ListAccount :many
Select * From account
Order By id
Limit $1
Offset $2;

-- name: UpdateAccount :one
Update account 
set balance =$2
where id =$1
returning *;


-- name: AddAccountBalance :one
Update account 
set balance = balance+ sqlc.arg(amount)
where id = sqlc.arg(id)
returning *;

-- name: DeleteAccount :exec
delete from account where id= $1;
