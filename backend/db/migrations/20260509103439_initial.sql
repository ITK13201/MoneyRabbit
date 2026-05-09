-- Create "categories" table
CREATE TABLE `categories` (
  `id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL,
  `icon` varchar(255) NOT NULL,
  `type` enum('income','expense','both') NOT NULL,
  `sort_order` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "category_rules" table
CREATE TABLE `category_rules` (
  `id` char(36) NOT NULL,
  `keyword` varchar(255) NOT NULL,
  `priority` bigint NOT NULL DEFAULT 0,
  `category_rules` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `category_rules_categories_rules` (`category_rules`),
  CONSTRAINT `category_rules_categories_rules` FOREIGN KEY (`category_rules`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "transactions" table
CREATE TABLE `transactions` (
  `id` char(36) NOT NULL,
  `date` timestamp NOT NULL,
  `description` varchar(255) NOT NULL,
  `amount` bigint NOT NULL,
  `import_format_id` enum('smbc_bank','smbc_card') NOT NULL,
  `imported_at` timestamp NOT NULL,
  `category_transactions` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `transactions_categories_transactions` (`category_transactions`),
  CONSTRAINT `transactions_categories_transactions` FOREIGN KEY (`category_transactions`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
