инициализация бд
создаем бд и роли, даем права
psql -U postgres -h 127.0.0.1 -p 6789 < timing/db/init.sql

создаем таблицы
psql -U postgres -h 127.0.0.1 -p 6789 -d timing_db < timing/db/tables.sql