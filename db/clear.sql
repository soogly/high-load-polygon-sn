
-- DO $$DECLARE count int;
-- BEGIN
-- SELECT count(*) INTO count FROM pg_roles WHERE rolname = 'go_app';
-- IF count > 0 THEN
--     EXECUTE 'DROP OWNED BY go_app CASCADE;';
-- END IF;
-- END$$;

DROP DATABASE IF EXISTS go_app_db;
DROP USER IF EXISTS go_app;


