CREATE TABLE IF NOT EXISTS users (
  id            SERIAL,
  email         varchar(255) NOT NULL,
  firstname     varchar(255) NOT NULL,
  lastname      varchar(255) NOT NULL,
  password      VARCHAR(255) NOT NULL
);
CREATE INDEX IF NOT EXISTS last_first_name_idx  ON users(lastname, firstname);

CREATE TABLE IF NOT EXISTS sessions (
  sessid        varchar(255) PRIMARY KEY,
  user_id       INTEGER NOT NULL,
  expires       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 days'
);
CREATE INDEX IF NOT EXISTS userid_idx  ON sessions(user_id);


GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO go_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO go_app;

INSERT INTO users (email, firstname, lastname, password) VALUES
('emma@mail.ru', 'Emma', 'Austen', 'password'),
('hz@mail.ru', 'Toto', 'Paolo', 'password'),
('lol@pol', 'Filipp', 'Fillipi4', 'password');