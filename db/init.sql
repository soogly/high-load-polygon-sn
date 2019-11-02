
-- CREATE USER go_app WITH PASSWORD 'go_pass';
-- CREATE USER 'go_app'@'%' IDENTIFIED BY 'go_pass';
CREATE DATABASE IF NOT EXISTS go_app_db;

GRANT ALL PRIVILEGES ON go_app_db.* TO 'go_app'@'%' IDENTIFIED BY 'go_pass';
GRANT INSERT ON go_app_db.* TO 'go_app'@'%' ;

-- CREATE INDEX firast_last_name_idx ON users (firstname(3), lastname(4));

