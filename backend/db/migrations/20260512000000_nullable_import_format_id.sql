-- +goose Up
ALTER TABLE `transactions`
  MODIFY COLUMN `import_format_id` enum('smbc_bank','smbc_card') NULL;

-- +goose Down
ALTER TABLE `transactions`
  MODIFY COLUMN `import_format_id` enum('smbc_bank','smbc_card') NOT NULL;
