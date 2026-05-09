-- +goose Up
-- Add "サブスク" category
INSERT IGNORE INTO `categories` (`id`, `name`, `color`, `icon`, `type`, `sort_order`) VALUES
  ('e7f8a9b0-c1d2-3456-efab-567890123456', 'サブスク', '#6366f1', '📺', 'expense', 85);

-- Add keyword rules for common subscription services
-- `category_rules` column is the FK to categories.id
INSERT IGNORE INTO `category_rules` (`id`, `keyword`, `priority`, `category_rules`) VALUES
  ('5b54e298-f5b7-436b-bd63-0b55d8bb64a4', 'NETFLIX',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('5106a3b4-f097-4da5-95e8-7eb1476ccc1f', 'SPOTIFY',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('f042f3c9-2b62-42b4-baa5-f2c8e340fc23', 'AMAZON PRIME',   50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('79a8ef38-8b1f-4651-bbd1-4809a805d1ae', 'APPLE.COM/BILL', 50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('be8e66bc-1782-4636-90c6-cebc6e44ca7f', 'YOUTUBE',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('ee5977b9-c06f-43bf-ac3a-db2f2f7bdfdd', 'DISNEY',         50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('a44d02cd-424c-4bc9-b1aa-a51d3f6d1607', 'HULU',           50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('9e2b04aa-46fd-41cf-9a5d-fee11c96ffc1', 'DAZN',           50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('d229a595-ef3d-4570-9645-c3ae08c0c161', 'ADOBE',          50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('b74e3da5-9413-46de-b19b-fb12a0fa5878', 'MICROSOFT',      50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('cbd8c9a3-86b0-439d-b983-b29d24e5ecb9', 'DROPBOX',        50, 'e7f8a9b0-c1d2-3456-efab-567890123456'),
  ('5ae15de7-ccf5-445f-96ff-2f906559fa60', 'NINTENDO',       50, 'e7f8a9b0-c1d2-3456-efab-567890123456');

-- +goose Down
DELETE FROM `category_rules` WHERE `id` IN (
  '5b54e298-f5b7-436b-bd63-0b55d8bb64a4',
  '5106a3b4-f097-4da5-95e8-7eb1476ccc1f',
  'f042f3c9-2b62-42b4-baa5-f2c8e340fc23',
  '79a8ef38-8b1f-4651-bbd1-4809a805d1ae',
  'be8e66bc-1782-4636-90c6-cebc6e44ca7f',
  'ee5977b9-c06f-43bf-ac3a-db2f2f7bdfdd',
  'a44d02cd-424c-4bc9-b1aa-a51d3f6d1607',
  '9e2b04aa-46fd-41cf-9a5d-fee11c96ffc1',
  'd229a595-ef3d-4570-9645-c3ae08c0c161',
  'b74e3da5-9413-46de-b19b-fb12a0fa5878',
  'cbd8c9a3-86b0-439d-b983-b29d24e5ecb9',
  '5ae15de7-ccf5-445f-96ff-2f906559fa60'
);
DELETE FROM `categories` WHERE `id` = 'e7f8a9b0-c1d2-3456-efab-567890123456';
