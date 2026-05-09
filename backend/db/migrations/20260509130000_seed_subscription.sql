-- +goose Up
-- Add "サブスク" category
INSERT IGNORE INTO `categories` (`id`, `name`, `color`, `icon`, `type`, `sort_order`) VALUES
  ('e7f8a9b0-c1d2-3456-efab-567890123456', 'サブスク', '#6366f1', '📺', 'expense', 85);

-- Add keyword rules for common subscription services
-- `category_rules` column is the FK to categories.id
INSERT IGNORE INTO `category_rules` (`id`, `keyword`, `priority`, `category_rules`) VALUES
  ('sub00001-0000-0000-0000-000000000001', 'NETFLIX',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000002', 'SPOTIFY',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000003', 'AMAZON PRIME',   50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000004', 'APPLE.COM/BILL', 50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000005', 'YOUTUBE',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000006', 'DISNEY',         50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000007', 'HULU',           50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000008', 'DAZN',           50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000009', 'ADOBE',          50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000010', 'MICROSOFT',      50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000011', 'DROPBOX',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('sub00001-0000-0000-0000-000000000012', 'NINTENDO',       50, 'e7f8a9b0-c1d2-3456-efab-567890123456');

-- +goose Down
DELETE FROM `category_rules` WHERE `id` LIKE 'sub00001-%';
DELETE FROM `categories` WHERE `id` = 'e7f8a9b0-c1d2-3456-efab-567890123456';
