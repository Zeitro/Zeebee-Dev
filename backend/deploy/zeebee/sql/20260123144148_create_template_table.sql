-- create base_templates table
CREATE TABLE "public"."base_templates" (
     "id" serial,
     "created_at" timestamptz NOT NULL DEFAULT now(),
     "updated_at" timestamptz NOT NULL DEFAULT now(),
     "is_del" int2 NOT NULL DEFAULT 0,
     "name" varchar(100) NOT NULL DEFAULT '',
     "type" varchar(20) NOT NULL DEFAULT '',
     "category" varchar(20) NOT NULL DEFAULT '',
     "tag" varchar(100) NOT NULL DEFAULT '',
     "content" text NOT NULL DEFAULT '',
     "variable_group_schema" jsonb NOT NULL DEFAULT '{}',
     "status" varchar(20) NOT NULL DEFAULT 'Active',
     "description" varchar(255) NOT NULL DEFAULT '',
     PRIMARY KEY ("id")
);

CREATE INDEX "idx_base_tpl_name" ON "public"."base_templates" ("name");
CREATE INDEX "idx_base_tpl_tag" ON "public"."base_templates" ("tag");
COMMENT ON COLUMN "public"."base_templates"."type" IS 'Borrower/Lender';
COMMENT ON COLUMN "public"."base_templates"."category" IS 'Email/Contract';
COMMENT ON COLUMN "public"."base_templates"."status" IS 'Active/Disabled';

-- create organization_templates table
CREATE TABLE "public"."organization_templates" (
   "id" serial,
   "created_at" timestamptz NOT NULL DEFAULT now(),
   "updated_at" timestamptz NOT NULL DEFAULT now(),
   "is_del" int2 NOT NULL DEFAULT 0,
   "organization_id" int4 NOT NULL DEFAULT 0,
   "name" varchar(100) NOT NULL DEFAULT '',
   "type" varchar(20) NOT NULL DEFAULT '',
   "category" varchar(20) NOT NULL DEFAULT '',
   "tag" varchar(100) NOT NULL DEFAULT '',
   "content" text NOT NULL DEFAULT '',
   "variable_group_schema" jsonb NOT NULL DEFAULT '{}',
   "status" varchar(20) NOT NULL DEFAULT 'Active',
   "description" varchar(255) NOT NULL DEFAULT '',
   PRIMARY KEY ("id")
);

CREATE INDEX "idx_org_tpl_name" ON "public"."organization_templates" ("name");
CREATE INDEX "idx_org_tpl_tag" ON "public"."organization_templates" ("tag");
CREATE INDEX "idx_org_tpl_org_id" ON "public"."organization_templates" ("organization_id");
COMMENT ON COLUMN "public"."organization_templates"."type" IS 'Borrower/Lender';
COMMENT ON COLUMN "public"."organization_templates"."category" IS 'Email/Contract';
COMMENT ON COLUMN "public"."organization_templates"."status" IS 'Active/Disabled';

ALTER TABLE "public"."organization_templates"
    ADD COLUMN "base_template_id" int4 NOT NULL DEFAULT 0;
CREATE INDEX "idx_org_base_tpl_id" ON "public"."organization_templates" ("base_template_id");