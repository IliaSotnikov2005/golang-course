-- name: CreateSubscription :one
INSERT INTO subscriptions (owner, repo)
VALUES ($1, $2)
RETURNING *;

-- name: ListSubscriptions :many
SELECT * FROM subscriptions
ORDER BY id;

-- name: DeleteSubscription :exec
DELETE FROM subscriptions
WHERE owner = $1 AND repo = $2;
