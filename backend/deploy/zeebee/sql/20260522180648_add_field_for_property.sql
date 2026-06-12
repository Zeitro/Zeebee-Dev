ALTER TABLE "public"."borrower_properties"
    ADD COLUMN "lot_area_sqft" decimal(8,2) NOT NULL DEFAULT 0;

ALTER TABLE "public"."loan_properties"
    ADD COLUMN "lot_area_sqft" decimal(8,2) NOT NULL DEFAULT 0;