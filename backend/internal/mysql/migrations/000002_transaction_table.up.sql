-- Transaction Table
CREATE TABLE transactions (
    id int unsigned AUTO_INCREMENT PRIMARY KEY,
    user_id int unsigned NOT NULL,
    order_id int unsigned NOT NULL,
    payment_method varchar(124) NOT NULL COMMENT "MPESA or STRIPE",
    amount decimal(10,2) NOT NULL,
    status boolean NOT NULL DEFAULT false,
    payment_details json NOT NULL,
    result_description text NOT NULL,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX transactions_index_0 ON transactions (id);
CREATE INDEX transactions_index_1 ON transactions (status);

-- Foreign Keys
-- ALTER TABLE transactions ADD FOREIGN KEY (user_id) REFERENCES users (id);
-- ALTER TABLE transactions ADD FOREIGN KEY (order_id) REFERENCES orders (id);

ALTER TABLE transactions ADD CONSTRAINT fk_transactions_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE transactions ADD CONSTRAINT fk_transactions_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE;