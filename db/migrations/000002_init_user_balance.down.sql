DROP INDEX IF EXISTS user_balance_user_id_currency;

DROP TABLE IF EXISTS user_balance;

ALTER TABLE user_balance DROP CONSTRAINT IF EXISTS user_balance_user_id_currency_unique;
