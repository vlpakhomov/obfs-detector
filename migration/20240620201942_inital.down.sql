BEGIN;

DROP TRIGGER IF EXISTS set_transaction_updated_at ON transaction;
DROP TABLE IF EXISTS transaction;
DROP TYPE IF EXISTS TRANSACTION_STATUS;

DROP TABLE IF EXISTS address;

DROP TRIGGER IF EXISTS set_scanner_updated_at ON scanner;
DROP TABLE IF EXISTS scanner;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS currency;

COMMIT;