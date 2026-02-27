ALTER TABLE "public"."apis"
    ADD COLUMN "auth_required" bool NOT NULL DEFAULT false,
  ADD COLUMN "permission_required" bool NOT NULL DEFAULT false;

ALTER TABLE "public"."resources"
    ADD COLUMN "enable_data_scope" bool NOT NULL DEFAULT false;