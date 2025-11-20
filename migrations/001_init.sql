-- Таблица команд
CREATE TABLE IF NOT EXISTS teams (
    name VARCHAR(100) PRIMARY KEY
);

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    team_name VARCHAR(100) REFERENCES teams(name),
    is_active BOOLEAN DEFAULT TRUE
);

-- Таблица Pull Request
CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(50) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(50) REFERENCES users(user_id),
    status VARCHAR(10) NOT NULL, -- OPEN|MERGED
    assigned_reviewers JSONB DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP NULL
);
