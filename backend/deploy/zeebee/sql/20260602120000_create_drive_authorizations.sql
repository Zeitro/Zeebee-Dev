-- drive_authorizations: one row per (organization_id, organization_type).
-- Uniqueness is enforced at the application layer via Upsert (mirrors email_authorizations).
CREATE TABLE "public"."drive_authorizations"
(
    "id"                    serial,
    "organization_id"       int4         NOT NULL DEFAULT 0,
    "organization_type"     varchar(20)  NOT NULL DEFAULT '',
    "auth_email"            varchar(255) NOT NULL DEFAULT '',
    "is_authorized"         bool         NOT NULL DEFAULT false,
    "drive_id"              varchar(255) NOT NULL DEFAULT '',
    "drive_name"            varchar(255) NOT NULL DEFAULT 'My Drive',
    "folder_name"           varchar(255) NOT NULL DEFAULT 'Zeebee',
    "root_folder_id"        varchar(255) NOT NULL DEFAULT '',
    "root_folder_url"       varchar(255) NOT NULL DEFAULT '',
    "token"                 jsonb        NOT NULL DEFAULT '{}'::jsonb,
    "authorized_by_user_id" int4         NOT NULL DEFAULT 0,
    "created_at"            timestamptz  NOT NULL DEFAULT now(),
    "updated_at"            timestamptz  NOT NULL DEFAULT now(),
    "is_del"                int2         NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_drive_authorizations_org" ON "public"."drive_authorizations" ("organization_id", "organization_type", "is_del");

COMMENT ON COLUMN "public"."drive_authorizations"."drive_id" IS 'empty = My Drive, otherwise Shared Drive ID';
COMMENT ON COLUMN "public"."drive_authorizations"."token" IS 'OAuth2 token (includes refresh token)';
