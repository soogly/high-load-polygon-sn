CREATE USER timing_ms WITH PASSWORD 'go_pass';
CREATE DATABASE timing_db WITH OWNER timing_ms;

GRANT ALL PRIVILEGES ON DATABASE timing_db TO timing_ms;