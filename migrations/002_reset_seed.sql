-- 时光酿:重置数据 + 扩到 12 坛酒(演示用)
-- 执行方式:
--   psql -U godemo_user -d shiguang_niang -h 127.0.0.1 -f /tmp/002_reset_seed.sql

BEGIN;

-- ============ 1. 清空业务数据(保留系列/酒窖等元数据) ============
DELETE FROM payments;
DELETE FROM claims;
DELETE FROM jar_metrics;
DELETE FROM jar_timeline;

-- 重置用户的默认 claim
UPDATE users SET default_claim_id = 0;

-- 删旧酒坛(后面重新插)
DELETE FROM wine_jars;

-- 重置序列(让 ID 从小开始)
ALTER SEQUENCE wine_jars_id_seq RESTART WITH 1;
ALTER SEQUENCE claims_id_seq RESTART WITH 1;
ALTER SEQUENCE payments_id_seq RESTART WITH 1;
ALTER SEQUENCE jar_metrics_id_seq RESTART WITH 1;
ALTER SEQUENCE jar_timeline_id_seq RESTART WITH 1;

-- ============ 2. 插 12 坛酒(全部 idle 可认领) ============
-- 系列 1=原酿 / 3=酒酿酒,酒窖 1=四坪 / 2=云岭 / 3=终南
INSERT INTO wine_jars (id, code, series_id, cellar_id, year, status, created_at) VALUES
  (1,  'BQ-0827', 1, 1, 2024, 'idle', NOW()),
  (2,  'BQ-0901', 2, 2, 2024, 'idle', NOW()),
  (3,  'BQ-1024', 3, 3, 2024, 'idle', NOW()),
  (4,  'BQ-0312', 1, 1, 2025, 'idle', NOW()),
  (5,  'BQ-0521', 1, 2, 2025, 'idle', NOW()),
  (6,  'BQ-0808', 2, 1, 2024, 'idle', NOW()),
  (7,  'BQ-0918', 2, 3, 2024, 'idle', NOW()),
  (8,  'BQ-1102', 2, 2, 2025, 'idle', NOW()),
  (9,  'BQ-1212', 3, 1, 2024, 'idle', NOW()),
  (10, 'BQ-0118', 3, 2, 2025, 'idle', NOW()),
  (11, 'BQ-0214', 1, 3, 2025, 'idle', NOW()),
  (12, 'BQ-0606', 2, 1, 2025, 'idle', NOW());

-- 让 BIGSERIAL 的 nextval 跳过已插入的固定 ID
SELECT setval('wine_jars_id_seq', 12);

COMMIT;

-- 验证
SELECT id, code, status, series_id, cellar_id FROM wine_jars ORDER BY id;
