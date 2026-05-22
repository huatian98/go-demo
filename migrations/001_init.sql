-- 时光酿 数据库初始化 (PostgreSQL 13+)
-- 执行方式:
--   sudo -u postgres psql -d shiguang_niang -f /path/to/001_init.sql

BEGIN;

-- ============ 用户 ============
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  openid VARCHAR(64) NOT NULL UNIQUE,
  unionid VARCHAR(64),
  nickname VARCHAR(64),
  avatar VARCHAR(255),
  phone VARCHAR(20),
  default_claim_id BIGINT DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE users IS '用户表';

-- ============ 酒品系列 ============
CREATE TABLE IF NOT EXISTS wine_series (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  description TEXT,
  cover_url VARCHAR(255),
  base_price DECIMAL(10,2) NOT NULL DEFAULT 0,
  sort INT NOT NULL DEFAULT 0,
  status SMALLINT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE wine_series IS '酒品系列';

-- ============ 酒窖 ============
CREATE TABLE IF NOT EXISTS cellars (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  address VARCHAR(255),
  province VARCHAR(32),
  city VARCHAR(32),
  capacity INT NOT NULL DEFAULT 0,
  available INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE cellars IS '酒窖';

-- ============ 单坛酒 ============
CREATE TABLE IF NOT EXISTS wine_jars (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(32) NOT NULL UNIQUE,
  series_id BIGINT NOT NULL,
  cellar_id BIGINT NOT NULL,
  year INT,
  cover_url VARCHAR(255),
  current_owner_id BIGINT DEFAULT 0,
  status VARCHAR(16) NOT NULL DEFAULT 'idle',
  claimed_at TIMESTAMPTZ,
  expected_ready_at TIMESTAMPTZ,
  version INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_wine_jars_series ON wine_jars (series_id);
CREATE INDEX IF NOT EXISTS idx_wine_jars_status ON wine_jars (status);
COMMENT ON TABLE wine_jars IS '单坛酒';

-- ============ 认领记录 ============
CREATE TABLE IF NOT EXISTS claims (
  id BIGSERIAL PRIMARY KEY,
  claim_no VARCHAR(32) NOT NULL UNIQUE,
  user_id BIGINT NOT NULL,
  jar_id BIGINT NOT NULL,
  cellar_id BIGINT NOT NULL,
  applicant_name VARCHAR(32) NOT NULL,
  contact_phone VARCHAR(32) NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_claims_user ON claims (user_id);
CREATE INDEX IF NOT EXISTS idx_claims_status ON claims (status);
COMMENT ON TABLE claims IS '认领记录';

-- ============ 支付订单 ============
CREATE TABLE IF NOT EXISTS payments (
  id BIGSERIAL PRIMARY KEY,
  claim_id BIGINT NOT NULL,
  out_trade_no VARCHAR(64) NOT NULL UNIQUE,
  channel VARCHAR(16) NOT NULL,
  amount DECIMAL(10,2) NOT NULL,
  status VARCHAR(16) NOT NULL DEFAULT 'pending',
  transaction_id VARCHAR(64),
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_payments_claim ON payments (claim_id);
COMMENT ON TABLE payments IS '支付订单';

-- ============ 酒坛实时指标 ============
CREATE TABLE IF NOT EXISTS jar_metrics (
  id BIGSERIAL PRIMARY KEY,
  wine_jar_id VARCHAR(32) NOT NULL,             -- 关联 wine_jars.code
  wine_ph DECIMAL(4,2) NOT NULL,                 -- 酸碱度
  ph_status VARCHAR(16),                          -- 稳定/偏高
  in_cellar_temp DECIMAL(4,1) NOT NULL,           -- 酒窖内温度
  in_cellar_humidity DECIMAL(4,1) NOT NULL,       -- 酒窖内湿度
  out_cellar_temp DECIMAL(4,1),                   -- 酒窖外温度
  out_cellar_humidity DECIMAL(4,1),               -- 酒窖外湿度
  breathing_state VARCHAR(32),                    -- 呼吸状态
  ai_narrative TEXT,                              -- AI 醒酒师文案
  recorded_at TIMESTAMPTZ NOT NULL,               -- 数据采集时间
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_jar_metrics_jar_time ON jar_metrics (wine_jar_id, recorded_at DESC);
COMMENT ON TABLE jar_metrics IS '酒坛实时指标(时序,wine_jar_id 关联 wine_jars.code)';

-- ============ 成长故事时间线 ============
CREATE TABLE IF NOT EXISTS jar_timeline (
  id BIGSERIAL PRIMARY KEY,
  jar_id BIGINT NOT NULL,
  event_type VARCHAR(32),
  title VARCHAR(64) NOT NULL,
  description TEXT,
  image_url VARCHAR(255),
  happened_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_jar_timeline_jar_time ON jar_timeline (jar_id, happened_at);
COMMENT ON TABLE jar_timeline IS '酒坛成长故事';

-- ============ 成分科普 ============
CREATE TABLE IF NOT EXISTS wine_components (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(32) NOT NULL,
  description TEXT,
  icon_url VARCHAR(255),
  sort INT NOT NULL DEFAULT 0
);
COMMENT ON TABLE wine_components IS '黄酒成分科普';

-- ============ 古法工艺 ============
CREATE TABLE IF NOT EXISTS craft_steps (
  id BIGSERIAL PRIMARY KEY,
  step_no INT NOT NULL,
  name VARCHAR(32) NOT NULL,
  description TEXT,
  image_url VARCHAR(255)
);
COMMENT ON TABLE craft_steps IS '古法酿造工艺';

-- ====================================================
-- 种子数据(演示用)
-- ====================================================

-- 系列
INSERT INTO wine_series (id, name, description, cover_url, base_price, sort, status)
VALUES
  (1, '十摊7春分系列', '春分时节古法酿造,选用屏南古田红曲,清雅醇厚', '/static/series/spring.jpg', 1299.00, 1, 1),
  (2, '十摊9秋分系列', '秋分时节窖藏,口感更醇厚饱满', '/static/series/autumn.jpg', 1499.00, 2, 1),
  (3, '十摊10冬至系列', '冬至深窖,适合长期陈酿', '/static/series/winter.jpg', 1899.00, 3, 1)
ON CONFLICT (id) DO NOTHING;

-- 酒窖
INSERT INTO cellars (id, name, address, province, city, capacity, available)
VALUES
  (1, '四平村古窖藏', '福建省宁德市屏南县四平村', '福建省', '宁德市', 500, 380),
  (2, '云岭古窖', '云南省大理州云龙县', '云南省', '大理州', 300, 250),
  (3, '终南山藏', '陕西省西安市长安区终南山', '陕西省', '西安市', 400, 320)
ON CONFLICT (id) DO NOTHING;

-- 酒坛(3 坛全部 idle 状态,可被认领)
INSERT INTO wine_jars (id, code, series_id, cellar_id, year, cover_url, status)
VALUES
  (1, 'BQ-0827', 1, 1, 2024, '/static/jars/bq-0827.png', 'idle'),
  (2, 'BQ-0901', 2, 2, 2024, '/static/jars/bq-0901.png', 'idle'),
  (3, 'BQ-1024', 3, 3, 2024, '/static/jars/bq-1024.png', 'idle')
ON CONFLICT (id) DO NOTHING;

-- 古法工艺
INSERT INTO craft_steps (step_no, name, description) VALUES
  (1, '浸米淘米', '糯米浸泡 24 小时,反复淘洗去杂质'),
  (2, '蒸饭摊凉', '木甑蒸饭半小时,出甑摊凉至 30°C'),
  (3, '拌曲下缸', '红曲、麦曲混匀拌入,装坛封口'),
  (4, '前发酵', '发酵 7-10 天,菌群活跃产酒'),
  (5, '入窖陈酿', '搬入古窖,慢呼吸 180 天以上'),
  (6, '过滤装坛', '压榨过滤后入小坛密封,等候开坛')
ON CONFLICT DO NOTHING;

-- 成分科普
INSERT INTO wine_components (name, description, sort) VALUES
  ('氨基酸', '黄酒含有 18 种氨基酸,其中 8 种是人体必需氨基酸', 1),
  ('多酚类', '抗氧化活性物质,有助于延缓衰老', 2),
  ('低聚糖', '促进肠道益生菌繁殖,改善消化', 3),
  ('麦角甾醇', '红曲特有,可调节胆固醇代谢', 4)
ON CONFLICT DO NOTHING;

-- 让 BIGSERIAL 的 nextval 跳过已插入的固定 ID
SELECT setval('wine_series_id_seq', GREATEST((SELECT MAX(id) FROM wine_series), 1));
SELECT setval('cellars_id_seq',    GREATEST((SELECT MAX(id) FROM cellars), 1));
SELECT setval('wine_jars_id_seq',  GREATEST((SELECT MAX(id) FROM wine_jars), 1));

COMMIT;
