CREATE TABLE IF NOT EXISTS `users` (
  `id` CHAR(36) UNIQUE NOT NULL,
  `email` VARCHAR(255) UNIQUE NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `subscription` BOOLEAN NOT NULL DEFAULT false COMMENT 'subscription to our blog posts',
  `role` VARCHAR(124) NOT NULL COMMENT 'USER or ADMIN',
  `refresh_token` VARCHAR(255) NOT NULL,
  `updated_by` CHAR(36) NOT NULL DEFAULT NULL,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  -- CONSTRAINT fk_updated_by FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `products` (
  `id` CHAR(36) UNIQUE NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) NOT NULL,
  `regular_price` DECIMAL(10,2) NOT NULL,
  `discounted_price` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `quantity` INT UNSIGNED NOT NULL DEFAULT 0,
  `category_id` INT UNSIGNED NOT NULL,
  `size_option` JSON NOT NULL DEFAULT ('{}'),
  `color_option` JSON NOT NULL DEFAULT ('{}'),
  `rating` float NOT NULL DEFAULT 0 COMMENT 'will be updated anytime a review is added',
  `seasonal` BOOLEAN NOT NULL DEFAULT false,
  `featured` BOOLEAN NOT NULL DEFAULT false,
  `img_urls` JSON NOT NULL,
  `updated_by` CHAR(36) NOT NULL DEFAULT null,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  -- CONSTRAINT fk_category_id FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  -- CONSTRAINT fk_updated_by FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `categories` (
  `id` INT PRIMARY KEY NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` TEXT NOT NULL DEFAULT 'No description'
);

CREATE TABLE IF NOT EXISTS `reviews` (
  `id` INT UNSIGNED PRIMARY KEY NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `product_id` CHAR(36) NOT NULL,
  `rating` INT UNSIGNED NOT NULL,
  `review` TEXT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  -- CONSTRAINT fk_user_id FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  -- CONSTRAINT fk_product_id FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `cart` (
  `user_id` CHAR(36) NOT NULL,
  `product_id` CHAR(36) NOT NULL,
  `quantity` INT UNSIGNED NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'will be used to check how long the cart has stayed'
  -- CONSTRAINT fk_user_id FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  -- CONSTRAINT fk_product_id FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `orders` (
  `id` CHAR(36) UNIQUE NOT NULL,
  `user_id` CHAR(36) NOT NULL,
  `amount` DECIMAL(10,2) NOT NULL COMMENT 'total amount of money for the order',
  `shipping_amount` DECIMAL(10,2) NOT NULL COMMENT 'shipping cost',
  `status` VARCHAR(255) NOT NULL DEFAULT 'PENDING' COMMENT 'PENDING, PROCESSING, SHIPPED or DELIVERED',
  `shipping_address` TEXT NOT NULL,
  `updated_by` CHAR(36) NOT NULL DEFAULT null,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  -- CONSTRAINT fk_user_id FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  -- CONSTRAINT fk_updated_by FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `order_items` (
  `order_id` CHAR(36) NOT NULL,
  `product_id` CHAR(36) NOT NULL,
  `quantity` INT UNSIGNED NOT NULL,
  `price` DECIMAL(10,2) NOT NULL,
  `color` VARCHAR(255) NOT NULL DEFAULT 'No color',
  `size` VARCHAR(255) NOT NULL DEFAULT 'No size'
  -- CONSTRAINT fk_order_id FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  -- CONSTRAINT fk_product_id FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `blogs` (
  `id` INT UNSIGNED PRIMARY KEY NOT NULL,
  `author` CHAR(36) NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `content` TEXT NOT NULL,
  `img_urls` JSON NOT NULL DEFAULT ('{}'),
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  -- CONSTRAINT fk_author FOREIGN KEY (`author`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX `users_index_0` ON `users` (`id`);
CREATE INDEX `users_index_1` ON `users` (`email`);
CREATE INDEX `products_index_2` ON `products` (`id`);
CREATE INDEX `reviews_index_3` ON `reviews` (`user_id`);
CREATE INDEX `reviews_index_4` ON `reviews` (`product_id`);
CREATE INDEX `cart_items_index_5` ON `cart_items` (`user_id`);
CREATE INDEX `cart_items_index_6` ON `cart_items` (`product_id`);
CREATE INDEX `orders_index_7` ON `orders` (`id`);
CREATE INDEX `orders_index_8` ON `orders` (`user_id`);
CREATE INDEX `orders_index_9` ON `orders` (`status`);

ALTER TABLE `products` ADD FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`);
ALTER TABLE `products` ADD FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`);
ALTER TABLE `reviews` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
ALTER TABLE `reviews` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);
ALTER TABLE `cart_items` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
ALTER TABLE `cart_items` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);
ALTER TABLE `orders` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
ALTER TABLE `orders` ADD FOREIGN KEY (`updated_by`) REFERENCES `users` (`id`);
ALTER TABLE `order_items` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);
ALTER TABLE `order_items` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);
ALTER TABLE `blogs` ADD FOREIGN KEY (`author`) REFERENCES `users` (`id`);
ALTER TABLE `users` ADD FOREIGN KEY (`users`) REFERENCES `users` (`id`);
