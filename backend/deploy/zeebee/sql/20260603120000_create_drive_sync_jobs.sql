-- drive_sync_jobs: one row per asynchronous "sync transaction to Google Drive" run.
-- The HTTP request only kicks off the job (status=Running) and returns immediately;
-- a background goroutine performs the actual Drive sync and updates this row when
-- done. The frontend polls the latest job for a transaction to show progress/result.
CREATE TABLE "public"."drive_sync_jobs"
(
    "id"                serial,
    "organization_id"   int4         NOT NULL DEFAULT 0,
    "organization_type" varchar(20)  NOT NULL DEFAULT '',
    "transaction_id"    int4         NOT NULL DEFAULT 0,
    "status"            varchar(20)  NOT NULL DEFAULT 'Running',
    "total"             int4         NOT NULL DEFAULT 0,
    "created"           int4         NOT NULL DEFAULT 0,
    "replaced"          int4         NOT NULL DEFAULT 0,
    "failed"            int4         NOT NULL DEFAULT 0,
    "folder_web_link"   varchar(512) NOT NULL DEFAULT '',
    "files"             jsonb        NOT NULL DEFAULT '[]'::jsonb,
    "error"             varchar(512) NOT NULL DEFAULT '',
    "created_at"        timestamptz  NOT NULL DEFAULT now(),
    "updated_at"        timestamptz  NOT NULL DEFAULT now(),
    "is_del"            int2         NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_drive_sync_jobs_tx" ON "public"."drive_sync_jobs" ("transaction_id", "organization_id", "is_del");

COMMENT ON COLUMN "public"."drive_sync_jobs"."status" IS 'Running,Success,PartialFailed,Failed';
COMMENT ON COLUMN "public"."drive_sync_jobs"."files" IS 'per-file sync result: [{file_name, drive_url, status}]';
