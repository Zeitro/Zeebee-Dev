CREATE TABLE "public"."role_resource_snapshots" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "organization_id" int4 NOT NULL default 0,
    "pricing_plan_id" int4 NOT NULL default 0,
    "role_id" int4 NOT NULL default 0,
    "resources" jsonb NOT NULL DEFAULT '[]'
);

CREATE INDEX "idx_role_resource_snapshots_organization_id" ON "public"."role_resource_snapshots" USING btree ("organization_id");
CREATE INDEX "idx_role_resource_snapshots_pricing_plan_id" ON "public"."role_resource_snapshots" USING btree ("pricing_plan_id");
CREATE INDEX "idx_role_resource_snapshots_role_id" ON "public"."role_resource_snapshots" USING btree ("role_id");
