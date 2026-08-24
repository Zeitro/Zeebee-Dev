ALTER TABLE "public"."project_management_task_assignees"
    ADD COLUMN "stage_id" int4 NOT NULL DEFAULT 0;

CREATE INDEX "idx_pm_assignee_stage"
    ON "public"."project_management_task_assignees" USING btree ("stage_id");
