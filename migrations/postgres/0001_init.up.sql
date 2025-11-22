-- Creating the relation for the models 'Team'.
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    team_name TEXT UNIQUE NOT NULL
);

-- Creating the relation for the models 'User'.
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    team_id INT NOT NULL REFERENCES teams(id) ON UPDATE CASCADE
);

-- Creating the relation for the models 'Pull request'.
CREATE TABLE IF NOT EXISTS pull_requests (
    id TEXT PRIMARY KEY,
    pr_name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status in ('OPEN', 'MERGED')),
    created_at TIMESTAMP,
    merged_at TIMESTAMP,
    author_id TEXT NOT NULL REFERENCES users(id) ON UPDATE CASCADE
);

-- Creating the relation connecting the models 'Pull request' and 'User'.
CREATE TABLE IF NOT EXISTS assigned_reviewers (
    pr_id TEXT NOT NULL REFERENCES pull_requests(id) ON UPDATE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON UPDATE CASCADE
);
