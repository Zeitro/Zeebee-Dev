-- Normalize organization.source to PascalCase (Manual | Search) to match Go constants.
UPDATE "public"."organizations"
SET source = 'Manual'
WHERE lower(source) = 'manual';

UPDATE "public"."organizations"
SET source = 'Search'
WHERE lower(source) = 'search';

ALTER TABLE "public"."organizations"
  ALTER COLUMN "source" SET DEFAULT 'Manual';

COMMENT ON COLUMN "public"."organizations"."source" IS 'Manual | Search';
