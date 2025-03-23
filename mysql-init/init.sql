CREATE DATABASE IF NOT EXISTS task_db;
CREATE USER IF NOT EXISTS 'task_user'@'%' IDENTIFIED WITH mysql_native_password BY 'task_password';
GRANT ALL PRIVILEGES ON task_db.* TO 'task_user'@'%';

CREATE DATABASE IF NOT EXISTS auth_db;
CREATE USER IF NOT EXISTS 'auth_user'@'%' IDENTIFIED WITH mysql_native_password BY 'auth_password';
GRANT ALL PRIVILEGES ON auth_db.* TO 'auth_user'@'%';

FLUSH PRIVILEGES;