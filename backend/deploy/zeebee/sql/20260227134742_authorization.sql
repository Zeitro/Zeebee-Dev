ALTER TABLE "public"."apis"
    ADD COLUMN "auth_required" bool NOT NULL DEFAULT false,
  ADD COLUMN "permission_required" bool NOT NULL DEFAULT false;

ALTER TABLE "public"."resources"
    ADD COLUMN "enable_data_scope" bool NOT NULL DEFAULT false;

ALTER TABLE "public"."resources"
    ALTER COLUMN "is_show" SET DEFAULT false;

DROP TABLE "public"."api_white_lists";

ALTER TABLE "public"."resources"
    ADD COLUMN "description" varchar(255) NOT NULL DEFAULT '';

ALTER TABLE "public"."resources"
    ADD COLUMN "is_public" bool NOT NULL DEFAULT false;

INSERT INTO "public"."apis" ("type", "name", "path", "method", "description", "status", "sort_order", "auth_required",
                             "permission_required")
VALUES ('Admin', 'sync data', '/api/v1/sync_data', 'GET', 'sync data', 'Active', 4, 'f', 'f');

CREATE TABLE "public"."casbin_rule"
(
    "id"    serial4,
    "ptype" varchar(100) COLLATE "pg_catalog"."default",
    "v0"    varchar(100) COLLATE "pg_catalog"."default",
    "v1"    varchar(100) COLLATE "pg_catalog"."default",
    "v2"    varchar(100) COLLATE "pg_catalog"."default",
    "v3"    varchar(100) COLLATE "pg_catalog"."default",
    "v4"    varchar(100) COLLATE "pg_catalog"."default",
    "v5"    varchar(100) COLLATE "pg_catalog"."default",
    CONSTRAINT "casbin_rule_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_casbin_rule" ON "public"."casbin_rule" USING btree (
    "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
    "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
    );
