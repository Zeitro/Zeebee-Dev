ALTER TABLE "public"."borrower_property_transactions"
    ADD COLUMN "title_holder" varchar(100) NOT NULL DEFAULT '',
  ADD COLUMN "contract_date" date,
  ADD COLUMN "emd_due_date" date,
  ADD COLUMN "contingency_removal_date" date,
  ADD COLUMN "tc" varchar(50) NOT NULL DEFAULT '',
  ADD COLUMN "cost_allocation" varchar(255) NOT NULL DEFAULT '';


CREATE TABLE "public"."property_transaction_escrows"
(
    "id"                               serial,
    "created_at"                       timestamptz  NOT NULL DEFAULT now(),
    "updated_at"                       timestamptz  NOT NULL DEFAULT now(),
    "is_del"                           int2         NOT NULL DEFAULT 0,
    "borrower_property_transaction_id" int4         NOT NULL DEFAULT 0,
    "no"                               varchar(50)  NOT NULL DEFAULT '',
    "escrow_company"                   varchar(100) NOT NULL DEFAULT '',
    "officer_name"                     varchar(30)  NOT NULL DEFAULT '',
    "officer_phone"                    varchar(20)  NOT NULL DEFAULT '',
    "officer_email"                    varchar(50)  NOT NULL DEFAULT '',
    "checklist"                        jsonb        NOT NULL DEFAULT '{}',
    PRIMARY KEY ("id")
);

CREATE INDEX "bpt_escrow_transaction_id" ON "public"."property_transaction_escrows" USING btree (
    "borrower_property_transaction_id"
    );

CREATE INDEX "bpt_escrow_no" ON "public"."property_transaction_escrows" USING btree (
    "no"
    );

CREATE TABLE "public"."property_transaction_properties"
(
    "id"                               serial,
    "created_at"                       timestamptz   NOT NULL DEFAULT now(),
    "updated_at"                       timestamptz   NOT NULL DEFAULT now(),
    "is_del"                           int2          NOT NULL DEFAULT 0,
    "borrower_property_transaction_id" int4          NOT NULL DEFAULT 0,
    "borrower_property_id"             int4          NOT NULL DEFAULT 0,
    "unit"                             int2          NOT NULL DEFAULT 1,
    "address"                          jsonb         NOT NULL DEFAULT '{}',
    "property_type"                    varchar(30)   NOT NULL DEFAULT '',
    "living_area_sqft"                 decimal(8, 2) NOT NULL DEFAULT 0,
    "bedrooms"                         int2          NOT NULL DEFAULT 0,
    "bathrooms"                        int2          NOT NULL DEFAULT 0,
    "lot_area_sqft"                    decimal(8, 2) NOT NULL DEFAULT 0,
    PRIMARY KEY ("id")
);

CREATE INDEX "bpt_transaction_id" ON "public"."property_transaction_properties" USING btree (
    "borrower_property_transaction_id"
    );
CREATE INDEX "bpt_property_id" ON "public"."property_transaction_properties" USING btree (
    "borrower_property_id"
    );

UPDATE "public"."document_templates"
SET "structure" = '{"common": [{"name": "Contract", "children": [{"name": "Original RPA", "sort_order": 1, "is_required": true,"file_parse_config":{"parser":"rpa"}}, {"name": "Assignment Agreement AOAA", "sort_order": 2, "is_required": false}, {"name": "Counter Offer", "sort_order": 3, "is_required": false}, {"name": "Addendum (Extension)", "sort_order": 4, "is_required": false}, {"name": "Seller financing Docs", "sort_order": 5, "is_required": false}, {"name": "JV Agreement", "sort_order": 6, "is_required": false}], "sort_order": 1, "is_required": true}, {"name": "Escrow", "children": [{"name": "Wire Instructions", "sort_order": 1, "is_required": true,"file_parse_config":{"parser":"escrow_info"}}, {"name": "EMD Payment Proof", "sort_order": 2, "is_required": true}, {"name": "EMD Receipt", "sort_order": 3, "is_required": true}, {"name": "Prelim Report", "sort_order": 4, "is_required": true}, {"name": "NHD Report", "sort_order": 5, "is_required": true}, {"name": "Colored map", "sort_order": 6, "is_required": false}, {"name": "FinCen", "sort_order": 7, "is_required": false}, {"name": "Buyer Open Package", "sort_order": 8, "is_required": false}], "sort_order": 2, "is_required": true}, {"name": "Disclosures", "sort_order": 3, "is_required": true}, {"name": "Loan", "children": [{"name": "Appraisal Report (Appraisal Payment confirmation and receipt)", "sort_order": 1, "is_required": true}, {"name": "LOI", "sort_order": 2, "is_required": true}, {"name": "Application Form", "sort_order": 3, "is_required": false}, {"name": "Loan Document", "sort_order": 4, "is_required": false}, {"name": "Contact form", "sort_order": 5, "is_required": false}], "sort_order": 4, "is_required": true}, {"name": "Insurance", "children": [{"name": "EOI", "sort_order": 1, "is_required": true}, {"name": "Payment Receipt", "sort_order": 2, "is_required": true}, {"name": "Full Policy", "sort_order": 3, "is_required": true}, {"name": "Flood Insurance", "sort_order": 4, "is_required": false}], "sort_order": 5, "is_required": true}, {"name": "Closing Package", "children": [{"name": "Final Buyer Statement", "sort_order": 1, "is_required": true}, {"name": "All Signed Documents", "sort_order": 2, "is_required": true}, {"name": "Recorded GD", "sort_order": 3, "is_required": true}, {"name": "Recorded DOT", "sort_order": 4, "is_required": false}, {"name": "Buyer Estimated Statement", "sort_order": 5, "is_required": false}], "sort_order": 6, "is_required": true}, {"name": "1031", "children": [{"name": "1031 Exchange Documents", "sort_order": 1, "is_required": false}], "sort_order": 7, "is_required": false}, {"name": "Lease", "children": [{"name": "Lease Agreement", "sort_order": 1, "is_required": false}, {"name": "TENANT ESTOPPEL CERTIFICATE", "sort_order": 2, "is_required": false}], "sort_order": 8, "is_required": false}, {"name": "HOA", "children": [{"name": "See \"HOA\" Sheet", "sort_order": 1, "is_required": false}], "sort_order": 9, "is_required": false}, {"name": "Solar", "children": [{"name": "See \"Solar\" Sheet", "sort_order": 1, "is_required": false}], "sort_order": 10, "is_required": false}]}'
WHERE "type" = 'BorrowerPropertyTransaction';