CREATE TABLE IF NOT EXISTS things (
  id            SERIAL,
  ttype         VARCHAR(20) NOT NULL,  -- events | regular | big dial | goal | shit
  title         VARCHAR(255) NOT NULL,
  creation_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  dude          INTEGER NOT NULL,
  comment       VARCHAR(655),
  priority      INTEGER NOT NULL DEFAULT 1000,
  duration      INTEGER DEFAULT 60,

  -- For `Events`
  when_it       TIMESTAMP DEFAULT NULL,

  -- For `regular events`
  step          INTEGER DEFAULT NULL,
  start_time    INTEGER ARRAY[2] DEFAULT '{-1,-1}'::int[],
  starts_from   DATE DEFAULT NULL,
  
  -- Days of week
  only_in       INTEGER ARRAY[7] DEFAULT '{-1,-1}'::int[],

  -- for `Shit` and `BigDeal`
  done          BOOLEAN DEFAULT FALSE,

  -- only  `Goal`
  big_deal      INTEGER DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS user_idx ON things (dude);


GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO timing_ms;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO timing_ms;

INSERT INTO things (ttype, title, dude, duration, step, start_time, starts_from) 
VALUES ('regular', 'Сон', 5, 420, 1, '{0,0}', CURRENT_TIMESTAMP);

INSERT INTO things (ttype, title, dude, duration, step, start_time, starts_from) 
VALUES ('regular', 'Пора спать', 5, 60, 1, '{23,0}', CURRENT_TIMESTAMP);

INSERT INTO things (ttype, title, dude, step, start_time, starts_from) 
VALUES ('regular', 'Завтрак', 5, 1, '{8,0}', CURRENT_TIMESTAMP);

INSERT INTO things (ttype, title, dude, step, start_time, starts_from) 
VALUES ('regular', 'Обед', 5, 1, '{12,30}', CURRENT_TIMESTAMP);

INSERT INTO things (ttype, title, dude, duration, step, start_time, starts_from) 
VALUES ('regular', 'Полдник', 5, 30, 1, '{16,0}', CURRENT_TIMESTAMP);

INSERT INTO things (ttype, title, dude, step, start_time, starts_from) 
VALUES ('regular', 'Ужин', 5, 1, '{19,0}', CURRENT_TIMESTAMP);