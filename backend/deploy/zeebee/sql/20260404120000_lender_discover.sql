-- Lender discovery via OpenAI: settings (singleton), run logs, per-item audit rows

CREATE TABLE "public"."lender_discover_settings"
(
    "id"                 serial4 PRIMARY KEY,
    "created_at"         timestamptz NOT NULL DEFAULT now(),
    "updated_at"         timestamptz NOT NULL DEFAULT now(),
    "is_del"             int2         NOT NULL DEFAULT 0,
    "schedule_enabled"   bool         NOT NULL DEFAULT false,
    "interval_days"      int4         NOT NULL DEFAULT 7,
    "scope_state"        varchar(8)   NOT NULL DEFAULT 'CA',
    "scope_extra"        text         NOT NULL DEFAULT '',
    "user_prompt"        text         NOT NULL DEFAULT '',
    "openai_model"       varchar(64)  NOT NULL DEFAULT 'gpt-4o-mini',
    "updated_by_user_id" int8         NOT NULL DEFAULT 0
);

CREATE TABLE "public"."lender_discover_runs"
(
    "id"              serial4 PRIMARY KEY,
    "created_at"      timestamptz NOT NULL DEFAULT now(),
    "updated_at"      timestamptz NOT NULL DEFAULT now(),
    "is_del"          int2         NOT NULL DEFAULT 0,
    "trigger_type"    varchar(32)  NOT NULL DEFAULT '',
    "status"          varchar(32)  NOT NULL DEFAULT '',
    "started_at"      timestamptz NOT NULL DEFAULT now(),
    "finished_at"     timestamptz,
    "error_message"   text,
    "stats"           jsonb        NOT NULL DEFAULT '{}',
    "prompt_snapshot" text         NOT NULL DEFAULT '',
    "scope_snapshot"  jsonb        NOT NULL DEFAULT '{}'
);

CREATE INDEX "idx_lender_discover_runs_started" ON "public"."lender_discover_runs" ("started_at" DESC);

CREATE TABLE "public"."lender_discover_run_items"
(
    "id"               serial4 PRIMARY KEY,
    "created_at"       timestamptz NOT NULL DEFAULT now(),
    "run_id"           int4         NOT NULL,
    "raw_name"         text         NOT NULL DEFAULT '',
    "normalized_key"   text         NOT NULL DEFAULT '',
    "action"           varchar(32)  NOT NULL DEFAULT '',
    "organization_id"  int8,
    "error_detail"     text         NOT NULL DEFAULT '',
    CONSTRAINT "fk_lender_discover_run_items_run" FOREIGN KEY ("run_id") REFERENCES "public"."lender_discover_runs" ("id") ON DELETE CASCADE
);

CREATE INDEX "idx_lender_discover_run_items_run_id" ON "public"."lender_discover_run_items" ("run_id");

INSERT INTO "public"."lender_discover_settings" ("schedule_enabled", "interval_days", "scope_state", "scope_extra", "user_prompt",
                                                 "openai_model")
SELECT false,
       7,
       'CA',
       'private mortgage lenders, hard money lenders, and private lending companies operating in this state',
       'Return only factual-sounding entries; leave unknown fields empty. Output must follow the requested JSON shape.',
       'gpt-4o-mini'
WHERE NOT EXISTS (SELECT 1 FROM "public"."lender_discover_settings" LIMIT 1);
