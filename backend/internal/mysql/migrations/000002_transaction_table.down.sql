ALTER TABLE transactions DROP FOREIGN KEY fk_transactions_user_id;
ALTER TABLE transactions DROP FOREIGN KEY fk_transactions_order_id;

DROP TABLE IF EXISTS transactions;