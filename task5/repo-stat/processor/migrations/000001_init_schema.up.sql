CREATE TABLE repositories (
    full_name    TEXT PRIMARY KEY,
    description  TEXT,
    stargazers   INTEGER NOT NULL DEFAULT 0,
    forks        INTEGER NOT NULL DEFAULT 0,
    created_at   TIMESTAMP WITH TIME ZONE,
    html_url     TEXT,
    last_cached  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
