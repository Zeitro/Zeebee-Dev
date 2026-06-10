-- pricing_plans
CREATE TABLE "public"."pricing_plans"
(
    "id"          serial,
    "created_at"  timestamptz  NOT NULL DEFAULT now(),
    "updated_at"  timestamptz  NOT NULL DEFAULT now(),
    "is_del"      int2         NOT NULL DEFAULT 0,
    "name"        varchar(100) NOT NULL DEFAULT '',
    "subtitle"    varchar(255) NOT NULL DEFAULT '',
    "description" varchar(255) NOT NULL DEFAULT '',
    "status"      varchar(50)  NOT NULL DEFAULT 'Active',
    "trial_days"  int2         NOT NULL DEFAULT 0,
    "sort_order"  int2         NOT NULL DEFAULT 6,
    "is_public"   bool         NOT NULL DEFAULT false,
    "is_default"  bool         NOT NULL DEFAULT false "type"        varchar(50)  NOT NULL DEFAULT '',
    "plan_type"   varchar(50)  NOT NULL DEFAULT 'Standard',
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_pricing_plan_name" ON "public"."pricing_plans" USING btree (
    "name"
    );

COMMENT
ON COLUMN "public"."pricing_plans"."status" IS 'Active,Disabled';
COMMENT
ON COLUMN "public"."pricing_plans"."type" IS 'Borrower,Lender,Admin';
COMMENT
ON COLUMN "public"."无pricing_plans标题"."plan_type" IS 'Standard,Custom';

-- pricing_plan_policies
CREATE TABLE "public"."pricing_plan_policies"
(
    "id"              serial,
    "created_at"      timestamptz    NOT NULL DEFAULT now(),
    "updated_at"      timestamptz    NOT NULL DEFAULT now(),
    "pricing_plan_id" int4           NOT NULL DEFAULT 0,
    "original_amount" decimal(10, 2) NOT NULL DEFAULT 0,
    "actual_amount"   decimal(10, 2) NOT NULL DEFAULT 0,
    "type"            varchar(20)    NOT NULL DEFAULT 'Monthly',
    PRIMARY KEY ("id")
)
;

CREATE INDEX "id_pr_policy_plan_id" ON "public"."pricing_plan_policies" USING btree (
    "pricing_plan_id"
    );

COMMENT
ON COLUMN "public"."pricing_plan_policies"."type" IS 'Monthly,Quarterly,Yearly';

-- pricing_plan_and_resources
CREATE TABLE "public"."pricing_plan_and_resources"
(
    "id"              serial,
    "pricing_plan_id" int4 NOT NULL DEFAULT 0,
    "resource_id"     int4 NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
)
;

CREATE INDEX "idx_pricing_plan_resource_plan_id" ON "public"."pricing_plan_and_resources" USING btree (
    "pricing_plan_id"
    );

CREATE INDEX "idx_pricing_plan_resource_resource_id" ON "public"."pricing_plan_and_resources" USING btree (
    "resource_id"
    );

-- organization_subscriptions
CREATE TABLE "public"."organization_subscriptions"
(
    "id"                     serial,
    "created_at"             timestamptz    NOT NULL DEFAULT now(),
    "updated_at"             timestamptz    NOT NULL DEFAULT now(),
    "is_del"                 int2           NOT NULL DEFAULT 0,
    "code"                   varchar(50)    NOT NULL DEFAULT '',
    "organization_id"        int4           NOT NULL DEFAULT 0,
    "pricing_plan_id"        int4           NOT NULL DEFAULT 0,
    "pricing_plan_policy_id" int4           NOT NULL DEFAULT 0,
    "status"                 varchar(50)    NOT NULL DEFAULT 'Active',
    "expire_date"            date           NOT NULL,
    "original_amount"        decimal(10, 2) NOT NULL DEFAULT 0,
    "actual_amount"          decimal(10, 2) NOT NULL DEFAULT 0,
    "is_auto_renew"          bool           NOT NULL DEFAULT false,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_org_subscri_code" ON "public"."organization_subscriptions" USING btree (
    "code"
    );

CREATE INDEX "idx_org_subscr_org_id" ON "public"."organization_subscriptions" USING btree (
    "organization_id"
    );

CREATE INDEX "idx_org_subscr_plan_id" ON "public"."organization_subscriptions" USING btree (
    "pricing_plan_id"
    );

CREATE INDEX "idx_org_subscr_policy_id" ON "public"."organization_subscriptions" USING btree (
    "pricing_plan_policy_id"
    );

CREATE INDEX ON "public"."organization_subscriptions" (
    );

COMMENT
ON COLUMN "public"."organization_subscriptions"."status" IS 'Trailing,Active,Expiring,Expired,Suspended';