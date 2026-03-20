package test

import "database/sql"

func DropTablesDB(db *sql.DB) error {
	_, err := db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;`)
	return err
}

func TruncateAllTablesDB(db *sql.DB) error {
	_, err := db.Exec(`
DO $$
DECLARE
    stmt TEXT;
BEGIN
    SELECT 'TRUNCATE TABLE ' || string_agg(quote_ident(tablename), ', ') || ' RESTART IDENTITY CASCADE;'
    INTO stmt
    FROM pg_tables
    WHERE schemaname = 'public';

    IF stmt IS NOT NULL THEN
        EXECUTE stmt;
    END IF;
END $$;
`)
	return err
}
