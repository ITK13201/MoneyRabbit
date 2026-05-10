-- +goose Up
INSERT IGNORE INTO `categories` (`id`, `name`, `color`, `icon`, `type`, `sort_order`) VALUES
  ('3175f88d-eb82-4264-84e1-ed82a83c963e', '旅行・宿泊', '#0ea5e9', '🏨', 'expense', 82);

INSERT IGNORE INTO `category_rules` (`id`, `keyword`, `priority`, `category_rules`) VALUES
  ('aecfe80e-9322-4ca8-9e5f-515ede289630', 'ホテル',  50, '3175f88d-eb82-4264-84e1-ed82a83c963e'),
  ('d20aacc7-1e6c-4d0c-a7a8-69539a9150c0', 'HOTEL',   50, '3175f88d-eb82-4264-84e1-ed82a83c963e'),
  ('1837a89e-1b8b-49ca-98ed-bbe15aaaf912', '旅館',    50, '3175f88d-eb82-4264-84e1-ed82a83c963e'),
  ('cb56e340-d697-47e0-ac4a-fe0756b33a9d', 'INN',     50, '3175f88d-eb82-4264-84e1-ed82a83c963e'),
  ('b26c0c87-18af-4ae5-a7bc-f4dfee1a13c3', 'RESORT',  50, '3175f88d-eb82-4264-84e1-ed82a83c963e'),
  ('5be33c8a-758b-433a-b987-1dcfa7daf0d8', 'リゾート', 50, '3175f88d-eb82-4264-84e1-ed82a83c963e'),
  ('b5d7c13d-7e4d-4606-8512-6267634b8d02', '宿泊',    50, '3175f88d-eb82-4264-84e1-ed82a83c963e');

-- +goose Down
DELETE FROM `category_rules` WHERE `id` IN (
  'aecfe80e-9322-4ca8-9e5f-515ede289630',
  'd20aacc7-1e6c-4d0c-a7a8-69539a9150c0',
  '1837a89e-1b8b-49ca-98ed-bbe15aaaf912',
  'cb56e340-d697-47e0-ac4a-fe0756b33a9d',
  'b26c0c87-18af-4ae5-a7bc-f4dfee1a13c3',
  '5be33c8a-758b-433a-b987-1dcfa7daf0d8',
  'b5d7c13d-7e4d-4606-8512-6267634b8d02'
);
DELETE FROM `categories` WHERE `id` = '3175f88d-eb82-4264-84e1-ed82a83c963e';
