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
    "enabled" bool NOT NULL DEFAULT true,
    "is_default" bool NOT NULL DEFAULT false,
    "source_template_id" int4 NOT NULL DEFAULT 0,
    "source_template_name" varchar(100) NOT NULL DEFAULT ''
);

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

CREATE TABLE "public"."project_management_workflow_transitions" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "workflow_id" int4 NOT NULL DEFAULT 0,
    "from_stage_id" int4 NOT NULL DEFAULT 0,
    "to_stage_id" int4 NOT NULL DEFAULT 0,
    "allowed_actor_types" text[] NOT NULL DEFAULT '{}'
);

CREATE TABLE "public"."project_management_projects" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "organization_id" int4 NOT NULL DEFAULT 0,
    "property_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "description" varchar(500) NOT NULL DEFAULT '',
    "cover_url" varchar(255) NOT NULL DEFAULT '',
    "owner_id" int4 NOT NULL DEFAULT 0
);

CREATE TABLE "public"."project_management_project_members" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "project_id" int4 NOT NULL DEFAULT 0,
    "organization_user_id" int4 NOT NULL DEFAULT 0,
    "role" varchar(32) NOT NULL DEFAULT 'TaskExecutor'
);

CREATE TABLE "public"."project_management_task_types" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "organization_id" int4 NOT NULL DEFAULT 0,
    "name" varchar(100) NOT NULL DEFAULT '',
    "description" varchar(500) NOT NULL DEFAULT '',
    "color" varchar(32) NOT NULL DEFAULT '',
    "enabled" bool NOT NULL DEFAULT true
);

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

CREATE TABLE "public"."project_management_task_assignees" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "task_id" int4 NOT NULL DEFAULT 0,
    "organization_user_id" int4 NOT NULL DEFAULT 0,
    "organization_user_name" varchar(100) NOT NULL DEFAULT ''
);

CREATE TABLE "public"."project_management_task_dependencies" (
    "id" serial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "is_del" int2 NOT NULL DEFAULT 0,
    "successor_task_id" int4 NOT NULL DEFAULT 0,
    "predecessor_task_id" int4 NOT NULL DEFAULT 0,
    "dependency_type" varchar(2) NOT NULL DEFAULT 'FS'
);

COMMENT ON COLUMN "public"."project_management_workflows"."kind" IS 'OrganizationTemplate or Project';
COMMENT ON COLUMN "public"."project_management_workflow_stages"."phase" IS 'Pending, InProgress or Finished';
COMMENT ON COLUMN "public"."project_management_workflow_transitions"."allowed_actor_types" IS 'ProjectOwner, ProjectAdmin, TaskCreator or TaskExecutor';
COMMENT ON COLUMN "public"."project_management_projects"."property_id" IS '0 means no Property binding';
COMMENT ON COLUMN "public"."project_management_project_members"."role" IS 'ProjectOwner, ProjectAdmin or TaskExecutor';
COMMENT ON COLUMN "public"."project_management_tasks"."parent_task_id" IS '0 means root task';
COMMENT ON COLUMN "public"."project_management_tasks"."task_type" IS 'immutable JSON snapshot: id, name, description and color';
COMMENT ON COLUMN "public"."project_management_task_dependencies"."dependency_type" IS 'FS, FF, SS or SF';

CREATE UNIQUE INDEX "uniq_project_management_default_workflow"
    ON "public"."project_management_workflows" ("organization_id")
    WHERE "kind" = 'OrganizationTemplate' AND "is_default" = true AND "is_del" = 0;
CREATE UNIQUE INDEX "uniq_project_management_project_workflow"
    ON "public"."project_management_workflows" ("project_id")
    WHERE "kind" = 'Project' AND "is_del" = 0;
CREATE UNIQUE INDEX "uniq_project_management_stage_sort_order"
    ON "public"."project_management_workflow_stages" ("workflow_id", "sort_order") WHERE "is_del" = 0;
CREATE UNIQUE INDEX "uniq_project_management_pending_stage"
    ON "public"."project_management_workflow_stages" ("workflow_id")
    WHERE "phase" = 'Pending' AND "is_del" = 0;
CREATE UNIQUE INDEX "uniq_project_management_transition"
    ON "public"."project_management_workflow_transitions" ("workflow_id", "from_stage_id", "to_stage_id") WHERE "is_del" = 0;
CREATE UNIQUE INDEX "uniq_project_management_member"
    ON "public"."project_management_project_members" ("project_id", "organization_user_id") WHERE "is_del" = 0;
CREATE UNIQUE INDEX "uniq_project_management_assignee"
    ON "public"."project_management_task_assignees" ("task_id", "organization_user_id") WHERE "is_del" = 0;
CREATE UNIQUE INDEX "uniq_project_management_dependency"
    ON "public"."project_management_task_dependencies" ("successor_task_id", "predecessor_task_id") WHERE "is_del" = 0;

CREATE INDEX "idx_project_management_project_property" ON "public"."project_management_projects" ("organization_id", "property_id", "is_del");
CREATE INDEX "idx_project_management_task_project_stage" ON "public"."project_management_tasks" ("project_id", "stage_id", "is_del");
CREATE INDEX "idx_project_management_task_parent" ON "public"."project_management_tasks" ("project_id", "parent_task_id", "is_del");
CREATE INDEX "idx_project_management_task_creator" ON "public"."project_management_tasks" ("organization_id", "creator_id", "is_del");
