CREATE TABLE user_balance (
	id SERIAL PRIMARY KEY,
	balance INT NOT NULL DEFAULT 0,
	currency VARCHAR(60) NOT NULL,
	user_id INT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

ALTER TABLE user_balance
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS user_balance_user_id_currency
	ON user_balance (user_id, currency);
ALTER TABLE user_balance ADD CONSTRAINT
	user_balance_user_id_currency_unique UNIQUE (user_id, currency);
