
CREATE TABLE IF NOT EXISTS users (
  id            INT(11) NOT NULL AUTO_INCREMENT,
  email         varchar(255) NOT NULL,
  firstname     varchar(255) NOT NULL,
  lastname      varchar(255) NOT NULL,
  password      VARCHAR(255) NOT NULL,
  CONSTRAINT pk PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS sessions (
  sessid        varchar(255),
  user_id       INTEGER NOT NULL,
  expires       DATETIME NOT NULL,
  CONSTRAINT pk PRIMARY KEY (sessid)
);


-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO go_app;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO go_app;

INSERT INTO users (email, firstname, lastname, password) VALUES
('emma@mail.ru', 'Emma', 'Austen', 'password'),
('hz@mail.ru', 'Toto', 'Paolo', 'password'),
('lol@pol', 'Filipp', 'Fillipi4', 'password');

