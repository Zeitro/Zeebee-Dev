-- 防止同一交易并发触发重复同步：用部分唯一索引在 DB 层硬保证
-- 「每个交易在同一组织下至多一个 Running 任务」。两个请求同时创建时，第二个会
-- 触发唯一冲突，应用层捕获后回退复用已存在的 Running 任务。

-- 1. 先清理历史遗留的重复 Running，每组只保留最新一条，
--    其余置为 Failed，否则下面的唯一索引会因已有重复而创建失败。
UPDATE "public"."drive_sync_jobs"
SET "status" = 'Failed',
    "error"  = 'superseded: duplicate running job cleaned up'
WHERE "status" = 'Running'
  AND "is_del" = 0
  AND "id" NOT IN (
      SELECT MAX("id")
      FROM "public"."drive_sync_jobs"
      WHERE "status" = 'Running' AND "is_del" = 0
      GROUP BY "transaction_id", "organization_id", "organization_type"
  );

-- 2. 部分唯一索引：仅约束未删除的 Running 行。
CREATE UNIQUE INDEX IF NOT EXISTS "uniq_drive_sync_jobs_running"
    ON "public"."drive_sync_jobs" ("transaction_id", "organization_id", "organization_type")
    WHERE "status" = 'Running' AND "is_del" = 0;
