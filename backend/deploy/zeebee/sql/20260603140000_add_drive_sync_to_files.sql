-- 为 files 增加 Google Drive 同步标记：
--   drive_url        同步后该文件在 Google Drive 中的可访问链接
--   drive_synced_at  最近一次成功同步到 Google Drive 的时间（NULL 表示尚未同步）
-- 前端依据 drive_synced_at 是否存在来显示 "Google Drive" 标签。
ALTER TABLE "public"."files"
    ADD COLUMN IF NOT EXISTS "drive_url" varchar(512) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS "drive_synced_at" timestamptz(6) NULL;

COMMENT ON COLUMN "public"."files"."drive_url" IS 'Google Drive 文件链接（同步后写入）';
COMMENT ON COLUMN "public"."files"."drive_synced_at" IS '最近一次同步到 Google Drive 的时间，NULL 表示未同步';
