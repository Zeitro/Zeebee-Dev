-- Organization display order: lower value = higher in list; 0 is min.
-- Default: discover (Not activated) = 1, active (manual create) = 2.
ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS "order" INT NOT NULL DEFAULT 2;

COMMENT ON COLUMN organizations."order" IS 'Display priority: 0 = top, 1 = discovered/not-activated, 2 = active default';

-- Back-fill existing rows based on source + status.
UPDATE organizations SET "order" = 1
WHERE source = 'Search' AND status = 'Not activated';
