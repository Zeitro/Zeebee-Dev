-- borrower_property_transactions
CREATE TABLE "public"."borrower_property_transactions"
(
    "id"                   serial,
    "borrower_property_id" int4           NOT NULL DEFAULT 0,
    "organization_id"      int4           NOT NULL DEFAULT 0,
    "code"                 varchar(20)    NOT NULL DEFAULT '',
    "type"                 varchar(20)    NOT NULL DEFAULT 'Purchase',
    "status"               varchar(20)    NOT NULL DEFAULT 'Pending',
    "coe_data"             date,
    "purchase_amount"      decimal(10, 2) NOT NULL DEFAULT 0,
    "emd_amount"           decimal(10, 2) NOT NULL DEFAULT 0,
    "property_type"        varchar(50)    NOT NULL DEFAULT '',
    "pm"                   varchar(255)   NOT NULL DEFAULT '',
    "business_strategy"    text[] COLLATE "pg_catalog"."default" DEFAULT '{}'::text[],
    "note"                 varchar(255)   NOT NULL DEFAULT '',
    "created_at"           timestamptz    NOT NULL DEFAULT now(),
    "updated_at"           timestamptz    NOT NULL DEFAULT now(),
    "is_del"               int2           NOT NULL DEFAULT 0,
    "checklist"            jsonb          NOT NULL DEFAULT '{}'::jsonb,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_bp_trans_property_id" ON "public"."borrower_property_transactions" ("borrower_property_id");
CREATE INDEX "idx_bp_trans_org_id" ON "public"."borrower_property_transactions" ("organization_id");
CREATE UNIQUE INDEX "idx_bp_trans_code" ON "public"."borrower_property_transactions" ("code");

COMMENT
ON COLUMN "public"."borrower_property_transactions"."code" IS 'unique code';
COMMENT
ON COLUMN "public"."borrower_property_transactions"."type" IS 'Purchase,Sale,Re-finance,Change Title,Assignment';
COMMENT
ON COLUMN "public"."borrower_property_transactions"."status" IS 'Pending,Closed,Canclled';
COMMENT
ON COLUMN "public"."borrower_property_transactions"."coe_data" IS 'close of end date';

-- operations
CREATE TABLE "public"."operations"
(
    "id"            serial,
    "entity_id"     int4         NOT NULL DEFAULT 0,
    "entity_type"   varchar(100) NOT NULL DEFAULT '',
    "operator_id"   int4         NOT NULL DEFAULT 0,
    "operator_name" varchar(50)  NOT NULL DEFAULT '',
    "operator_type" varchar(50)  NOT NULL DEFAULT '',
    "operation"     varchar(50)  NOT NULL DEFAULT '',
    "note"          varchar(255) NOT NULL DEFAULT '',
    "metadata"      jsonb        NOT NULL DEFAULT '{}',
    "created_at"    timestamptz  NOT NULL DEFAULT now(),
    "updated_at"    timestamptz  NOT NULL DEFAULT now(),
    "is_del"        int2         NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_operation_entity" ON "public"."operations" ("entity_id","entity_type");

ALTER TABLE "public"."operations"
    ADD COLUMN "organization_id" int4 NOT NULL DEFAULT 0;

CREATE INDEX "idx_operation_org_id" ON "public"."operations" ("organization_id");

-- entity_documents
CREATE TABLE "public"."entity_documents"
(
    "id"               serial,
    "organization_id" int4 NOT NULL DEFAULT 0,
    "entity_type"      varchar(50) NOT NULL DEFAULT '',
    "entity_id"        int4        NOT NULL DEFAULT 0,
    "root_document_id" int4        NOT NULL DEFAULT 0,
    "created_at"       timestamptz NOT NULL DEFAULT now(),
    "updated_at"       timestamptz NOT NULL DEFAULT now(),
    "is_del"           int2        NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_entity_doc_entity_id" ON "public"."entity_documents" ("entity_type","entity_id");
CREATE INDEX "idx_entity_doc_root_id" ON "public"."entity_documents" ("root_document_id");
CREATE INDEX "idx_entity_doc_org_id" ON "public"."entity_documents" ("organization_id");

-- documents
CREATE TABLE "public"."documents"
(
    "id"                   serial,
    "parent_id"            int4         NOT NULL DEFAULT 0,
    "entity_document_id"   int4         NOT NULL DEFAULT 0,
    "type"                 varchar(50)  NOT NULL DEFAULT '',
    "name"                 varchar(255) NOT NULL DEFAULT '',
    "level"                int4         NOT NULL DEFAULT 0,
    "sort_order"           int4         NOT NULL DEFAULT 9,
    "last_uploaded_at"     timestamptz,
    "metadata"             jsonb        NOT NULL DEFAULT '{}',
    "created_at"           timestamptz  NOT NULL DEFAULT now(),
    "updated_at"           timestamptz  NOT NULL DEFAULT now(),
    "is_del"               int2         NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_doc_pid" ON "public"."documents" ("parent_id");
CREATE INDEX "idx_doc_name" ON "public"."documents" ("name");
CREATE INDEX "idx_doc_entity_document_id" ON "public"."documents" ("entity_document_id");

-- files
ALTER TABLE "public"."files"
    ADD COLUMN "file_name" varchar(255) NOT NULL DEFAULT '',
  ADD COLUMN "storage_type" varchar(20) NOT NULL DEFAULT 'Aws S3',
  ADD COLUMN "external_id" varchar(255) NOT NULL DEFAULT '',
  ADD COLUMN "upload_type" varchar(20) NOT NULL DEFAULT 'Manual';

CREATE INDEX "idx_files_external_id" ON "public"."files" ("external_id");

COMMENT
ON COLUMN "public"."files"."storage_type" IS 'Aws S3,Google Drive ';
COMMENT
ON COLUMN "public"."files"."upload_type" IS 'Manual,Google Drive';

-- document_templates
CREATE TABLE "public"."document_templates"
(
    "id"         serial,
    "name"       varchar(255) NOT NULL DEFAULT '',
    "type"       varchar(20)  NOT NULL DEFAULT '',
    "structure"  jsonb         NOT NULL DEFAULT '{}',
    "status"     varchar(20)  NOT NULL DEFAULT '',
    "created_at" timestamptz  NOT NULL DEFAULT now(),
    "updated_at" timestamptz  NOT NULL DEFAULT now(),
    "is_del"     int2         NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

-- loans
ALTER TABLE public.loans
    ADD column borrower_property_transaction_id int4 DEFAULT 0 NOT NULL;
CREATE INDEX idx_loan_borrower_property_trans_id ON public.loans (borrower_property_transaction_id);