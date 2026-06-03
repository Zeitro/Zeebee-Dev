-- 同步时会把「我们之前上传、但系统侧已删除」的 Drive 文件移到回收站，
-- removed 记录本次同步移除的文件数，用于结果弹窗展示。
ALTER TABLE "public"."drive_sync_jobs"
    ADD COLUMN IF NOT EXISTS "removed" int4 NOT NULL DEFAULT 0;

COMMENT ON COLUMN "public"."drive_sync_jobs"."removed" IS '本次同步从 Drive 移除（回收站）的过期文件数';
