CREATE TABLE user_transactions (
	id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  amount INT NOT NULL,
  currency VARCHAR(60) NOT NULL,
  bank_account_number VARCHAR(30) NOT NULL,
  bank_name VARCHAR(30) NOT NULL,
  image_url TEXT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

ALTER TABLE user_transactions
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
