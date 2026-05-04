-- name: GetRepository :one
SELECT * FROM repositories
WHERE full_name = $1 LIMIT 1;

-- name: ListAllRepositories :many
SELECT * FROM repositories ORDER BY full_name;

-- name: UpsertRepository :exec
INSERT INTO repositories (
    full_name, description, stargazers, forks, created_at, html_url, last_cached
) VALUES (
    $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP
)
ON CONFLICT (full_name) DO UPDATE SET
    description = EXCLUDED.description,
    stargazers = EXCLUDED.stargazers,
    forks = EXCLUDED.forks,
    html_url = EXCLUDED.html_url,
    last_cached = CURRENT_TIMESTAMP;
