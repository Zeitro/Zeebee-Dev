-- stripe payment integration

-- organizations: one organization maps to one Stripe Customer
ALTER TABLE "public"."organizations"
    ADD COLUMN "stripe_customer_id" varchar(64) NOT NULL DEFAULT '';

CREATE INDEX "idx_org_stripe_customer_id" ON "public"."organizations" USING btree (
    "stripe_customer_id"
    );

-- pricing_plans: Stripe Product mapping (lazily created at checkout)
ALTER TABLE "public"."pricing_plans"
    ADD COLUMN "stripe_product_id" varchar(64) NOT NULL DEFAULT '';

-- pricing_plan_policies: Stripe Price mapping (lazily created at checkout, cleared on amount/type change)
ALTER TABLE "public"."pricing_plan_policies"
    ADD COLUMN "stripe_price_id" varchar(64) NOT NULL DEFAULT '';

-- organization_subscriptions: subscription source and Stripe subscription mapping
ALTER TABLE "public"."organization_subscriptions"
    ADD COLUMN "source" varchar(20) NOT NULL DEFAULT 'Manual',
    ADD COLUMN "stripe_subscription_id" varchar(64) NOT NULL DEFAULT '';

CREATE INDEX "idx_org_subscr_stripe_sub_id" ON "public"."organization_subscriptions" USING btree (
    "stripe_subscription_id"
    );

COMMENT
ON COLUMN "public"."organization_subscriptions"."source" IS 'Manual,Stripe';

-- stripe_events: webhook event idempotency ledger
CREATE TABLE "public"."stripe_events"
(
    "id"         serial,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del"     int2        NOT NULL DEFAULT 0,
    "event_id"   varchar(64) NOT NULL DEFAULT '',
    "event_type" varchar(64) NOT NULL DEFAULT '',
    "status"     varchar(20) NOT NULL DEFAULT 'Processing',
    "payload"    jsonb       NOT NULL DEFAULT '{}',
    "error"      text        NOT NULL DEFAULT '',
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_stripe_events_event_id" ON "public"."stripe_events" USING btree (
    "event_id"
    );

COMMENT
ON COLUMN "public"."stripe_events"."status" IS 'Processing,Processed,Failed,Skipped';

-- payment_transactions: payment records synced from Stripe invoices
CREATE TABLE "public"."payment_transactions"
(
    "id"                           serial,
    "created_at"                   timestamptz    NOT NULL DEFAULT now(),
    "updated_at"                   timestamptz    NOT NULL DEFAULT now(),
    "is_del"                       int2           NOT NULL DEFAULT 0,
    "organization_id"              int4           NOT NULL DEFAULT 0,
    "organization_subscription_id" int4           NOT NULL DEFAULT 0,
    "stripe_invoice_id"            varchar(64)    NOT NULL DEFAULT '',
    "stripe_subscription_id"       varchar(64)    NOT NULL DEFAULT '',
    "amount"                       decimal(10, 2) NOT NULL DEFAULT 0,
    "currency"                     varchar(10)    NOT NULL DEFAULT 'usd',
    "status"                       varchar(20)    NOT NULL DEFAULT '',
    "period_start"                 date,
    "period_end"                   date,
    "paid_at"                      timestamptz,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "idx_payment_tx_invoice_id" ON "public"."payment_transactions" USING btree (
    "stripe_invoice_id"
    );

CREATE INDEX "idx_payment_tx_org_id" ON "public"."payment_transactions" USING btree (
    "organization_id"
    );

COMMENT
ON COLUMN "public"."payment_transactions"."status" IS 'Paid,Failed,Refunded';

-- self-serve subscription APIs: login required, but exempt from the pricing
-- plan permission gate so owners of expired organizations can still pay.
-- sync_api only inserts missing rows, so these flags are not overwritten.
INSERT INTO "public"."apis" ("type", "name", "path", "method", "description", "status", "sort_order",
                             "auth_required", "permission_required")
VALUES ('Borrower', 'public pricing plans', '/api/v1/pricing_plans/public', 'GET', 'list public pricing plans', 'Active', 1, 't', 'f'),
       ('Borrower', 'current subscription', '/api/v1/organization/subscription/current', 'GET', 'current organization subscription', 'Active', 2, 't', 'f'),
       ('Borrower', 'subscription checkout', '/api/v1/organization/subscription/checkout', 'POST', 'create stripe checkout session', 'Active', 3, 't', 'f'),
       ('Borrower', 'cancel auto renewal', '/api/v1/organization/subscription/cancel_renewal', 'POST', 'cancel auto renewal', 'Active', 4, 't', 'f'),
       ('Borrower', 'resume auto renewal', '/api/v1/organization/subscription/resume_renewal', 'POST', 'resume auto renewal', 'Active', 5, 't', 'f'),
       ('Lender', 'public pricing plans', '/api/v1/pricing_plans/public', 'GET', 'list public pricing plans', 'Active', 1, 't', 'f'),
       ('Lender', 'current subscription', '/api/v1/organization/subscription/current', 'GET', 'current organization subscription', 'Active', 2, 't', 'f'),
       ('Lender', 'subscription checkout', '/api/v1/organization/subscription/checkout', 'POST', 'create stripe checkout session', 'Active', 3, 't', 'f'),
       ('Lender', 'cancel auto renewal', '/api/v1/organization/subscription/cancel_renewal', 'POST', 'cancel auto renewal', 'Active', 4, 't', 'f'),
       ('Lender', 'resume auto renewal', '/api/v1/organization/subscription/resume_renewal', 'POST', 'resume auto renewal', 'Active', 5, 't', 'f');
