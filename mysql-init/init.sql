CREATE DATABASE IF NOT EXISTS auth_db;
CREATE USER IF NOT EXISTS 'auth_user'@'%' IDENTIFIED WITH mysql_native_password BY 'auth_password';
GRANT ALL PRIVILEGES ON auth_db.* TO 'auth_user'@'%';

CREATE DATABASE IF NOT EXISTS task_db;
CREATE USER IF NOT EXISTS 'task_user'@'%' IDENTIFIED WITH mysql_native_password BY 'task_password';
GRANT ALL PRIVILEGES ON task_db.* TO 'task_user'@'%';

FLUSH PRIVILEGES;

-- Switch to the task_db database
USE task_db;

-- Drop the tasks table if it exists
DROP TABLE IF EXISTS tasks;

-- Create the tasks table
CREATE TABLE tasks (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255),
    content VARCHAR(4096) NOT NULL,
    stylized_content VARCHAR(8192),
    discarded BOOLEAN DEFAULT FALSE,
    status TINYINT(8) NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(36) NOT NULL,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    modified_by VARCHAR(36) NOT NULL
);

-- Create indexes for filtering and user tasks
CREATE INDEX idx_filtered_tasks ON tasks (discarded, status);
CREATE INDEX idx_user_tasks ON tasks (created_by);


-- Switch to the auth_db database
USE auth_db;

-- Drop the users table if it exists
DROP TABLE IF EXISTS users;

-- Create the tasks table
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    enc_password VARCHAR(64) NOT NULL,
    discarded BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for filtering and user tasks
CREATE INDEX idx_users ON users (username);

INSERT INTO users (id, username, enc_password)
VALUES ('user-id-abc', 'pavan', '$2a$12$83h4Mjus72N0cudx9upBGed9EXnCrPdd4sK2PG7WhxU36HlqehuYm');
INSERT INTO users (id, username, enc_password)
VALUES ('user-id-xyz', 'harshit', '$2a$12$npRO3ZOgQygBNXPoPbEIgOhu1vFBUhvyV3nnEqRmx8QzQLuUZNx0q');

