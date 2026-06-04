-- borrower_property_transactions
CREATE TABLE "public"."borrower_property_transactions"
(
    "id"                   serial,
    "borrower_property_id" int4           NOT NULL DEFAULT 0,
    "organization_id"      int4           NOT NULL DEFAULT 0,
    "code"                 varchar(20)    NOT NULL DEFAULT '',
    "type"                 varchar(20)    NOT NULL DEFAULT 'Purchase',
    "status"               varchar(20)    NOT NULL DEFAULT 'Pending',
    "coe_date"             date,
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
ON COLUMN "public"."borrower_property_transactions"."status" IS 'Pending,Closed,Cancelled';
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
    "type" varchar(50) NOT NULL DEFAULT '',
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

-- migrate existing strategy varchar to business_strategy text[]
ALTER TABLE public.borrower_property_transactions
    RENAME COLUMN strategy TO business_strategy;

ALTER TABLE public.borrower_property_transactions
    ALTER COLUMN business_strategy TYPE text[] USING
        CASE WHEN business_strategy = '' THEN '{}'::text[]
             ELSE string_to_array(business_strategy, ',') END;

ALTER TABLE public.borrower_property_transactions
    ALTER COLUMN business_strategy SET DEFAULT '{}'::text[];

-- document_templates insert
INSERT INTO "public"."document_templates" ("name", "type", "structure", "status")
VALUES ('Borrower Property Transaction Documents', 'BorrowerPropertyTransaction', '{
  "common": [
    {
      "name": "Contract",
      "children": [
        {
          "name": "Original RPA",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "Assignment Agreement AOAA",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "Counter Offer",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Addendum (Extension)",
          "sort_order": 4,
          "is_required": true
        },
        {
          "name": "Seller financing Docs",
          "sort_order": 5,
          "is_required": true
        },
        {
          "name": "JV Agreement",
          "sort_order": 6,
          "is_required": true
        }
      ],
      "sort_order": 1,
      "is_required": true
    },
    {
      "name": "Escrow",
      "children": [
        {
          "name": "Wire Instructions",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "EMD Payment Proof",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "EMD Receipt",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Prelim Report",
          "sort_order": 4,
          "is_required": true
        },
        {
          "name": "NHD Report",
          "sort_order": 5,
          "is_required": true
        },
        {
          "name": "Colored map",
          "sort_order": 6,
          "is_required": true
        },
        {
          "name": "FinCen",
          "sort_order": 7,
          "is_required": true
        },
        {
          "name": "Buyer Open Package",
          "sort_order": 8,
          "is_required": true
        },
        {
          "name": "LOI",
          "sort_order": 9,
          "is_required": true
        }
      ],
      "sort_order": 2,
      "is_required": true
    },
    {
      "name": "Disclosures",
      "sort_order": 3,
      "is_required": true
    },
    {
      "name": "Loan",
      "children": [
        {
          "name": "Application Form",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "Loan Document",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "Contact form",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Appraisal Report (Appraisal Payment confirmation and receipt)",
          "sort_order": 4,
          "is_required": true
        }
      ],
      "sort_order": 4,
      "is_required": true
    },
    {
      "name": "Insurance",
      "children": [
        {
          "name": "EOI",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "Payment Receipt",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "Full Policy",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Flood Insurance",
          "sort_order": 4,
          "is_required": false
        }
      ],
      "sort_order": 5,
      "is_required": true
    },
    {
      "name": "Closing Package",
      "children": [
        {
          "name": "Buyer Estimated Statement",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "Final Buyer Statement",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "All Signed Documents",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Recorded GD",
          "sort_order": 4,
          "is_required": true
        },
        {
          "name": "Recorded DOT",
          "sort_order": 5,
          "is_required": true
        }
      ],
      "sort_order": 6,
      "is_required": true
    },
    {
      "name": "1031",
      "children": [
        {
          "name": "1031 Exchange Documents",
          "sort_order": 1,
          "is_required": false
        }
      ],
      "sort_order": 7,
      "is_required": false
    },
    {
      "name": "Lease",
      "children": [
        {
          "name": "Lease Agreement",
          "sort_order": 1,
          "is_required": false
        },
        {
          "name": "TENANT ESTOPPEL CERTIFICATE",
          "sort_order": 2,
          "is_required": false
        }
      ],
      "sort_order": 8,
      "is_required": false
    },
    {
      "name": "HOA",
      "children": [
        {
          "name": "See \"HOA\" Sheet",
          "sort_order": 1,
          "is_required": false
        }
      ],
      "sort_order": 9,
      "is_required": false
    },
    {
      "name": "Solar",
      "children": [
        {
          "name": "See \"Solar\" Sheet",
          "sort_order": 1,
          "is_required": false
        }
      ],
      "sort_order": 10,
      "is_required": false
    }
  ]
}', 'Active');

update "public"."document_templates"
set "structure" = '{
  "common": [
    {
      "name": "Contract",
      "children": [
        {
          "name": "Original RPA",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "Assignment Agreement AOAA",
          "sort_order": 2,
          "is_required": false
        },
        {
          "name": "Counter Offer",
          "sort_order": 3,
          "is_required": false
        },
        {
          "name": "Addendum (Extension)",
          "sort_order": 4,
          "is_required": false
        },
        {
          "name": "Seller financing Docs",
          "sort_order": 5,
          "is_required": false
        },
        {
          "name": "JV Agreement",
          "sort_order": 6,
          "is_required": false
        }
      ],
      "sort_order": 1,
      "is_required": true
    },
    {
      "name": "Escrow",
      "children": [
        {
          "name": "Wire Instructions",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "EMD Payment Proof",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "EMD Receipt",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Prelim Report",
          "sort_order": 4,
          "is_required": true
        },
        {
          "name": "NHD Report",
          "sort_order": 5,
          "is_required": true
        },
        {
          "name": "Colored map",
          "sort_order": 6,
          "is_required": false
        },
        {
          "name": "FinCen",
          "sort_order": 7,
          "is_required": false
        },
        {
          "name": "Buyer Open Package",
          "sort_order": 8,
          "is_required": false
        }
      ],
      "sort_order": 2,
      "is_required": true
    },
    {
      "name": "Disclosures",
      "sort_order": 3,
      "is_required": true
    },
    {
      "name": "Loan",
      "children": [
        {
          "name": "Appraisal Report (Appraisal Payment confirmation and receipt)",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "LOI",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "Application Form",
          "sort_order": 3,
          "is_required": false
        },
        {
          "name": "Loan Document",
          "sort_order": 4,
          "is_required": false
        },
        {
          "name": "Contact form",
          "sort_order": 5,
          "is_required": false
        }
      ],
      "sort_order": 4,
      "is_required": true
    },
    {
      "name": "Insurance",
      "children": [
        {
          "name": "EOI",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "Payment Receipt",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "Full Policy",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Flood Insurance",
          "sort_order": 4,
          "is_required": false
        }
      ],
      "sort_order": 5,
      "is_required": true
    },
    {
      "name": "Closing Package",
      "children": [
        {
          "name": "Final Buyer Statement",
          "sort_order": 1,
          "is_required": true
        },
        {
          "name": "All Signed Documents",
          "sort_order": 2,
          "is_required": true
        },
        {
          "name": "Recorded GD",
          "sort_order": 3,
          "is_required": true
        },
        {
          "name": "Recorded DOT",
          "sort_order": 4,
          "is_required": false
        },
        {
          "name": "Buyer Estimated Statement",
          "sort_order": 5,
          "is_required": false
        }
      ],
      "sort_order": 6,
      "is_required": true
    },
    {
      "name": "1031",
      "children": [
        {
          "name": "1031 Exchange Documents",
          "sort_order": 1,
          "is_required": false
        }
      ],
      "sort_order": 7,
      "is_required": false
    },
    {
      "name": "Lease",
      "children": [
        {
          "name": "Lease Agreement",
          "sort_order": 1,
          "is_required": false
        },
        {
          "name": "TENANT ESTOPPEL CERTIFICATE",
          "sort_order": 2,
          "is_required": false
        }
      ],
      "sort_order": 8,
      "is_required": false
    },
    {
      "name": "HOA",
      "children": [
        {
          "name": "See \"HOA\" Sheet",
          "sort_order": 1,
          "is_required": false
        }
      ],
      "sort_order": 9,
      "is_required": false
    },
    {
      "name": "Solar",
      "children": [
        {
          "name": "See \"Solar\" Sheet",
          "sort_order": 1,
          "is_required": false
        }
      ],
      "sort_order": 10,
      "is_required": false
    }
  ]
}'
where "type" = 'BorrowerPropertyTransaction';

ALTER TABLE "public"."borrower_properties"
    ADD COLUMN "status" text[] NOT NULL DEFAULT '{}'::text[];

CREATE INDEX "idx_property_status" ON "public"."borrower_properties" USING gin ("status");