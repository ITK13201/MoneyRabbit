-- +goose Up
INSERT IGNORE INTO `categories` (`id`, `name`, `color`, `icon`, `type`, `sort_order`) VALUES
  ('be3882b0-4e06-4322-85f5-59c164de0e3e', '投資',   '#f59e0b', '📈', 'expense', 115),
  ('27fd50a2-8a30-47d7-9928-bd9d7054ea84', '積立',   '#0d9488', '🏦', 'expense', 116);

-- +goose Down
DELETE FROM `categories` WHERE `id` IN (
  'be3882b0-4e06-4322-85f5-59c164de0e3e',
  '27fd50a2-8a30-47d7-9928-bd9d7054ea84'
);
