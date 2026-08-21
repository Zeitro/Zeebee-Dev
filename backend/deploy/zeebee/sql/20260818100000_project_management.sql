-- Borrower project management tables and indexes.
-- The migration uses the repository's soft-delete convention: is_del = 0.

CREATE TABLE "public"."project_management_workflows" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "kind" varchar(32) NOT NULL DEFAULT 'OrganizationTemplate',
    "organization_id" int4 NOT NULL DEFAULT 0,
    "project_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "description" varchar(500) NOT NULL DEFAULT '',
    "is_default" bool NOT NULL DEFAULT false,
    "source_template_id" int4 NOT NULL DEFAULT 0,
    "source_template_name" varchar(100) NOT NULL DEFAULT ''
);

COMMENT ON COLUMN "public"."project_management_workflows"."kind" IS 'OrganizationTemplate or Project';

CREATE INDEX "idx_pm_org_id"
    ON "public"."project_management_workflows" using btree ("organization_id");
CREATE INDEX "idx_pm_proj_id"
    ON "public"."project_management_workflows" using btree ("project_id");

CREATE TABLE "public"."project_management_workflow_stages" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "workflow_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "phase" varchar(32) NOT NULL DEFAULT 'Pending',
    "sort_order" int4 NOT NULL DEFAULT 0
);

COMMENT ON COLUMN "public"."project_management_workflow_stages"."phase" IS 'Pending, InProgress or Finished';

CREATE INDEX "idx_pm_stage_sort_order"
    ON "public"."project_management_workflow_stages" USING btree ("workflow_id", "sort_order");
CREATE INDEX "idx_pm_stage_workflow_id"
    ON "public"."project_management_workflow_stages" USING btree ("workflow_id");

CREATE TABLE "public"."project_management_workflow_transitions" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "workflow_id" int4 NOT NULL DEFAULT 0,
    "from_stage_id" int4 NOT NULL DEFAULT 0,
    "to_stage_id" int4 NOT NULL DEFAULT 0,
    "allowed_actor_types" jsonb NOT NULL DEFAULT '[]'
);

COMMENT ON COLUMN "public"."project_management_workflow_transitions"."allowed_actor_types" IS 'ProjectOwner, ProjectAdmin, TaskCreator or TaskExecutor';

CREATE INDEX "idx_pm_transition"
    ON "public"."project_management_workflow_transitions" USING btree ("workflow_id", "from_stage_id", "to_stage_id");

CREATE TABLE "public"."project_management_projects" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "organization_id" int4 NOT NULL DEFAULT 0,
    "entity_id" int4 NOT NULL DEFAULT 0,
    "entity_type" varchar(100) NOT NULL DEFAULT '',
    "name" varchar(100) NOT NULL DEFAULT '',
    "description" varchar(500) NOT NULL DEFAULT '',
    "cover_url" varchar(255) NOT NULL DEFAULT '',
    "owner_id" int4 NOT NULL DEFAULT 0,
    "owner_name" varchar(100) NOT NULL DEFAULT ''
);

COMMENT ON COLUMN "public"."project_management_projects"."entity_id" IS 'ID of the bound entity';
COMMENT ON COLUMN "public"."project_management_projects"."entity_type" IS 'Type of the bound entity, e.g. Property or Transaction';

CREATE INDEX "idx_pm_project_entity"
    ON "public"."project_management_projects" USING btree ("organization_id", "entity_id", "entity_type");

CREATE TABLE "public"."project_management_project_members" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "project_id" int4 NOT NULL DEFAULT 0,
    "organization_user_id" int4 NOT NULL DEFAULT 0,
    "role" varchar(32) NOT NULL DEFAULT 'TaskExecutor'
);

COMMENT ON COLUMN "public"."project_management_project_members"."role" IS 'ProjectOwner, ProjectAdmin or TaskExecutor';

CREATE INDEX "idx_pm_member"
    ON "public"."project_management_project_members" USING btree ("project_id", "organization_user_id");

CREATE TABLE "public"."project_management_task_types" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "organization_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "description" varchar(500) NOT NULL DEFAULT '',
    "color" varchar(32) NOT NULL DEFAULT '',
    "sort_order" int4 NOT NULL DEFAULT 0
);

CREATE INDEX "idx_pm_task_type_organization_id"
    ON "public"."project_management_task_types" USING btree ("organization_id");

INSERT INTO "public"."project_management_task_types" ("organization_id", "name", "color", "sort_order")
VALUES
    (0, 'Task', '#1890FF', 1),
    (0, 'Defect', '#F5222D', 2),
    (0, 'Requirement', '#FAAD14', 3),
    (0, 'Improvement', '#52C41A', 4);

CREATE TABLE "public"."project_management_tasks" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "project_id" int4 NOT NULL DEFAULT 0,
    "organization_id" int4 NOT NULL DEFAULT 0,
    "parent_task_id" int4 NOT NULL DEFAULT 0,
    "stage_id" int4 NOT NULL DEFAULT 0,
    "creator_id" int4 NOT NULL DEFAULT 0,
    "creator_name" varchar(100) NOT NULL DEFAULT '',
    "title" varchar(255) NOT NULL DEFAULT '',
    "description" text NOT NULL DEFAULT '',
    "priority" varchar(32) NOT NULL DEFAULT '',
    "start_at" timestamptz NULL,
    "end_at" timestamptz NULL,
    "task_type" jsonb NOT NULL DEFAULT '{}'
);

COMMENT ON COLUMN "public"."project_management_tasks"."parent_task_id" IS '0 means root task';
COMMENT ON COLUMN "public"."project_management_tasks"."task_type" IS 'immutable JSON snapshot: id, name, description and color';

CREATE INDEX "idx_pm_task_project_stage" ON "public"."project_management_tasks" USING btree ("project_id", "stage_id");
CREATE INDEX "idx_pm_task_parent" ON "public"."project_management_tasks" USING btree ("project_id", "parent_task_id");
CREATE INDEX "idx_pm_task_creator" ON "public"."project_management_tasks" USING btree ("organization_id", "creator_id");

CREATE TABLE "public"."project_management_task_assignees" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "task_id" int4 NOT NULL DEFAULT 0,
    "organization_user_id" int4 NOT NULL DEFAULT 0,
    "organization_user_name" varchar(100) NOT NULL DEFAULT ''
);

CREATE INDEX "idx_pm_assignee"
    ON "public"."project_management_task_assignees" USING btree ("task_id", "organization_user_id","organization_user_name");

CREATE TABLE "public"."project_management_task_dependencies" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "task_id" int4 NOT NULL DEFAULT 0,
    "dependency_task_id" int4 NOT NULL DEFAULT 0,
    "dependency_type" varchar(16) NOT NULL DEFAULT 'PreTask',
    "dependency_rule" varchar(2) NOT NULL DEFAULT 'FS'
);

COMMENT ON COLUMN "public"."project_management_task_dependencies"."dependency_type" IS 'PreTask or PostTask';
COMMENT ON COLUMN "public"."project_management_task_dependencies"."dependency_rule" IS 'FS, FF, SS or SF';

CREATE INDEX "idx_pm_dependency_task"
    ON "public"."project_management_task_dependencies" USING btree ("task_id");
CREATE INDEX "idx_pm_dependency_related_task"
    ON "public"."project_management_task_dependencies" USING btree ("dependency_task_id");
