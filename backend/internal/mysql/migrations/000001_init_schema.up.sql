-- Users table
CREATE TABLE users (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  email varchar(255) UNIQUE NOT NULL,
  password varchar(255) NOT NULL,
  subscription boolean NOT NULL DEFAULT false COMMENT 'subscription to our blog posts',
  role varchar(124) NOT NULL COMMENT 'USER or ADMIN',
  refresh_token text NOT NULL,
  updated_by int unsigned,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- -- Categories table
CREATE TABLE categories (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  name varchar(255) NOT NULL,
  description text NOT NULL
);

-- -- Products table
CREATE TABLE products (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  name varchar(255) NOT NULL,
  description varchar(255) NOT NULL,
  regular_price decimal(10,2) NOT NULL,
  discounted_price decimal(10,2) NOT NULL DEFAULT 0.00,
  quantity int unsigned NOT NULL DEFAULT 0,
  category_id int unsigned NOT NULL,
  size_option json NOT NULL,
  color_option json NOT NULL,
  rating float NOT NULL DEFAULT 0 COMMENT 'will be updated anytime a review is added',
  seasonal boolean NOT NULL DEFAULT false,
  featured boolean NOT NULL DEFAULT false,
  img_urls json NOT NULL,
  updated_by int unsigned NOT NULL,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Reviews table
CREATE TABLE reviews (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  user_id int unsigned NOT NULL,
  product_id int unsigned NOT NULL,
  rating int unsigned NOT NULL,
  review text NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- -- Cart table
CREATE TABLE cart (
  user_id int unsigned NOT NULL,
  product_id int unsigned NOT NULL,
  quantity int unsigned NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'will be used to check how long the cart has stayed'
);

-- -- Orders table
CREATE TABLE orders (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  user_id int unsigned NOT NULL,
  amount decimal(10,2) NOT NULL COMMENT 'total amount of money for the order',
  shipping_amount decimal(10,2) NOT NULL COMMENT 'shipping cost',
  status varchar(255) NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING, PROCESSING, SHIPPED or DELIVERED',
  shipping_address text NOT NULL,
  updated_by int unsigned NOT NULL,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Order items table
CREATE TABLE order_items (
  order_id int unsigned NOT NULL,
  product_id int unsigned NOT NULL,
  quantity int unsigned NOT NULL,
  price decimal(10,2) NOT NULL,
  color varchar(255) NOT NULL DEFAULT 'No color',
  size varchar(255) NOT NULL DEFAULT 'No size'
);

-- Blogs table
CREATE TABLE blogs (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  author int unsigned NOT NULL,
  title varchar(255) NOT NULL,
  content text NOT NULL,
  img_urls json NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX users_index_0 ON users (id);
CREATE INDEX users_index_1 ON users (email);
CREATE INDEX products_index_2 ON products (id);
CREATE INDEX reviews_index_3 ON reviews (user_id);
CREATE INDEX reviews_index_4 ON reviews (product_id);
CREATE INDEX orders_index_7 ON orders (id);
CREATE INDEX orders_index_8 ON orders (user_id);
CREATE INDEX orders_index_9 ON orders (status);

-- Foreign Keys
-- ALTER TABLE users ADD FOREIGN KEY (updated_by) REFERENCES users (id);
-- ALTER TABLE products ADD FOREIGN KEY (category_id) REFERENCES categories (id);
-- ALTER TABLE products ADD FOREIGN KEY (updated_by) REFERENCES users (id);
-- ALTER TABLE reviews ADD FOREIGN KEY (user_id) REFERENCES users (id);
-- ALTER TABLE reviews ADD FOREIGN KEY (product_id) REFERENCES products (id);
-- ALTER TABLE cart ADD FOREIGN KEY (user_id) REFERENCES users (id);
-- ALTER TABLE cart ADD FOREIGN KEY (product_id) REFERENCES products (id);
-- ALTER TABLE orders ADD FOREIGN KEY (user_id) REFERENCES users (id);
-- ALTER TABLE orders ADD FOREIGN KEY (updated_by) REFERENCES users (id);
-- ALTER TABLE order_items ADD FOREIGN KEY (order_id) REFERENCES orders (id);
-- ALTER TABLE order_items ADD FOREIGN KEY (product_id) REFERENCES products (id);
-- ALTER TABLE blogs ADD FOREIGN KEY (author) REFERENCES users (id);

ALTER TABLE users ADD CONSTRAINT fk_users_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE products ADD CONSTRAINT fk_products_category_id FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE;
ALTER TABLE products ADD CONSTRAINT fk_products_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE reviews ADD CONSTRAINT fk_reviews_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE reviews ADD CONSTRAINT fk_reviews_product_id FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE;
ALTER TABLE cart ADD CONSTRAINT fk_cart_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE cart ADD CONSTRAINT fk_cart_product_id FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE;
ALTER TABLE orders ADD CONSTRAINT fk_orders_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE orders ADD CONSTRAINT fk_orders_updated_by FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE order_items ADD CONSTRAINT fk_order_items_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE;
ALTER TABLE order_items ADD CONSTRAINT fk_order_items_product_id FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE;
ALTER TABLE blogs ADD CONSTRAINT fk_blogs_author FOREIGN KEY (author) REFERENCES users (id) ON DELETE CASCADE