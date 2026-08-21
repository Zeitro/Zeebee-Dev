ALTER TABLE "public"."project_management_workflows"
    DROP COLUMN IF EXISTS "enabled";

WITH defaults AS (
    SELECT id, row_number() OVER (PARTITION BY organization_id ORDER BY id) AS sequence
    FROM "public"."project_management_workflows"
    WHERE kind = 'OrganizationTemplate' AND is_default = true AND is_del = 0
)
UPDATE "public"."project_management_workflows" AS workflows
SET is_default = false
FROM defaults
WHERE workflows.id = defaults.id AND defaults.sequence > 1;

CREATE UNIQUE INDEX "idx_pm_workflow_default"
    ON "public"."project_management_workflows" ("organization_id")
    WHERE kind = 'OrganizationTemplate' AND is_default = true AND is_del = 0;
