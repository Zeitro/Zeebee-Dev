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
