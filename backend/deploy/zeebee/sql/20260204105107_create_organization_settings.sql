ALTER TABLE "public"."base_templates"
    ADD COLUMN "is_configurable" bool NOT NULL DEFAULT false;

ALTER TABLE "public"."organization_templates"
    ADD COLUMN "has_configured" bool NOT NULL DEFAULT false;

-- create organization_settings table
CREATE TABLE "public"."organization_settings"
(
    "id"              serial4,
    "created_at"      timestamptz  NOT NULL DEFAULT now(),
    "updated_at"      timestamptz  NOT NULL DEFAULT now(),
    "is_del"          int2         NOT NULL DEFAULT 0,
    "organization_id" int4         NOT NULL DEFAULT 0,
    "purpose"         varchar(100) NOT NULL DEFAULT '',
    "enabled"         bool         NOT NULL DEFAULT false,
    "setting"         jsonb        NOT NULL DEFAULT '{}',
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_org_settings_org_id" ON "public"."organization_settings" ("organization_id");
CREATE INDEX "idx_org_settings_purpose" ON "public"."organization_settings" ("purpose");