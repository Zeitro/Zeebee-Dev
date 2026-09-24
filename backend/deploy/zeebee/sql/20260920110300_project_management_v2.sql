-- Borrower project management v1 -> v2 forward migration.
-- PostgreSQL-only. Keep every statement in one DO block because some SQL-file
-- runners check each statement out from a pool instead of pinning BEGIN..COMMIT.
-- A DO block is one server statement and therefore rolls back atomically on error.
DO $pm_v2_migration$
DECLARE
    bad_ids text;
BEGIN

-- Do not touch a v1 table until every target row has been materialized below.
CREATE TABLE pm_v2_projects_staging (LIKE project_management_projects INCLUDING ALL);
ALTER TABLE pm_v2_projects_staging ALTER COLUMN id DROP DEFAULT;
CREATE SEQUENCE pm_v2_projects_staging_id_seq OWNED BY pm_v2_projects_staging.id;
ALTER TABLE pm_v2_projects_staging ALTER COLUMN id SET DEFAULT nextval('pm_v2_projects_staging_id_seq');
ALTER TABLE pm_v2_projects_staging
    ADD COLUMN project_type varchar(32) NOT NULL DEFAULT 'Ordinary',
    ADD COLUMN start_at date,
    ADD COLUMN end_at date,
    ADD COLUMN period_value integer NOT NULL DEFAULT 0,
    ADD COLUMN period_unit varchar(16) NOT NULL DEFAULT '';

CREATE TABLE pm_v2_templates_staging (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(), is_del smallint NOT NULL DEFAULT 0,
    template_type varchar(32) NOT NULL, organization_id integer NOT NULL DEFAULT 0,
    name varchar(100) NOT NULL DEFAULT '', description varchar(500) NOT NULL DEFAULT '',
    status varchar(32) NOT NULL DEFAULT 'Active', is_default boolean NOT NULL DEFAULT false, sort_order integer NOT NULL DEFAULT 0,
    source_project_id integer NOT NULL DEFAULT 0
);

CREATE TABLE pm_v2_stages_staging (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(), is_del smallint NOT NULL DEFAULT 0,
    entity_type varchar(32) NOT NULL, entity_id integer NOT NULL, name varchar(100) NOT NULL DEFAULT '', description varchar(500) NOT NULL DEFAULT '',
    phase varchar(32) NOT NULL, color varchar(32) NOT NULL DEFAULT '', sort_order integer NOT NULL DEFAULT 0
);

CREATE TABLE pm_v2_workflows_staging (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(), is_del smallint NOT NULL DEFAULT 0,
    entity_type varchar(32) NOT NULL, entity_id integer NOT NULL, source_workflow_id integer NOT NULL DEFAULT 0,
    source_workflow_name varchar(100) NOT NULL DEFAULT '', name varchar(100) NOT NULL DEFAULT '', description varchar(500) NOT NULL DEFAULT '', is_default boolean NOT NULL DEFAULT false
);
CREATE TABLE pm_v2_workflow_stages_staging (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(), is_del smallint NOT NULL DEFAULT 0,
    workflow_id integer NOT NULL, stage_id integer NOT NULL
);
CREATE TABLE pm_v2_workflow_transitions_staging (
    id serial PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(), updated_at timestamptz NOT NULL DEFAULT now(), is_del smallint NOT NULL DEFAULT 0,
    workflow_id integer NOT NULL, name varchar(100) NOT NULL DEFAULT '', from_stage_id integer NOT NULL, to_stage_id integer NOT NULL,
    allowed_actor_types jsonb NOT NULL DEFAULT '[]'::jsonb
);

CREATE TABLE pm_v2_tasks_staging (LIKE project_management_tasks INCLUDING ALL);
ALTER TABLE pm_v2_tasks_staging ALTER COLUMN id DROP DEFAULT;
CREATE SEQUENCE pm_v2_tasks_staging_id_seq OWNED BY pm_v2_tasks_staging.id;
ALTER TABLE pm_v2_tasks_staging ALTER COLUMN id SET DEFAULT nextval('pm_v2_tasks_staging_id_seq');
ALTER TABLE pm_v2_tasks_staging RENAME COLUMN project_id TO entity_id;
ALTER TABLE pm_v2_tasks_staging ADD COLUMN entity_type varchar(32) NOT NULL DEFAULT 'Project';
ALTER TABLE pm_v2_tasks_staging ADD COLUMN workflow_id integer;
ALTER TABLE pm_v2_tasks_staging ALTER COLUMN stage_id DROP NOT NULL;
ALTER TABLE pm_v2_tasks_staging ALTER COLUMN phase DROP NOT NULL;
ALTER TABLE pm_v2_tasks_staging ALTER COLUMN stage_id DROP DEFAULT;
ALTER TABLE pm_v2_tasks_staging ALTER COLUMN phase DROP DEFAULT;
CREATE TABLE pm_v2_task_assignees_staging (LIKE project_management_task_assignees INCLUDING ALL);
ALTER TABLE pm_v2_task_assignees_staging ALTER COLUMN id DROP DEFAULT;
CREATE SEQUENCE pm_v2_task_assignees_staging_id_seq OWNED BY pm_v2_task_assignees_staging.id;
ALTER TABLE pm_v2_task_assignees_staging ALTER COLUMN id SET DEFAULT nextval('pm_v2_task_assignees_staging_id_seq');
ALTER TABLE pm_v2_task_assignees_staging ALTER COLUMN stage_id DROP NOT NULL;
ALTER TABLE pm_v2_task_assignees_staging ALTER COLUMN stage_id DROP DEFAULT;
CREATE TABLE pm_v2_task_dependencies_staging (LIKE project_management_task_dependencies INCLUDING ALL);
ALTER TABLE pm_v2_task_dependencies_staging ALTER COLUMN id DROP DEFAULT;
CREATE SEQUENCE pm_v2_task_dependencies_staging_id_seq OWNED BY pm_v2_task_dependencies_staging.id;
ALTER TABLE pm_v2_task_dependencies_staging ALTER COLUMN id SET DEFAULT nextval('pm_v2_task_dependencies_staging_id_seq');

CREATE TEMP TABLE pm_v2_workflow_map (old_id bigint PRIMARY KEY, new_id bigint NOT NULL, entity_type varchar(32) NOT NULL, entity_id bigint NOT NULL) ON COMMIT DROP;
CREATE TEMP TABLE pm_v2_source_workflow_map (old_source_template_id bigint PRIMARY KEY, new_source_workflow_id bigint NOT NULL) ON COMMIT DROP;
CREATE TEMP TABLE pm_v2_stage_map (old_id bigint PRIMARY KEY, new_id bigint NOT NULL) ON COMMIT DROP;


-- Preconditions name every unsafe v1 object. Dirty input is rejected, never merged.
    SELECT string_agg(id::text, ', ' ORDER BY id) INTO bad_ids
    FROM project_management_workflows WHERE is_del = 0 AND project_id <> 0
    GROUP BY project_id HAVING count(*) <> 1 LIMIT 1;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: project has zero or multiple active workflows: %', bad_ids; END IF;

    SELECT string_agg(w.id::text, ', ' ORDER BY w.id) INTO bad_ids
    FROM project_management_workflows w
    LEFT JOIN project_management_workflow_stages s ON s.workflow_id = w.id AND s.is_del = 0
    WHERE w.is_del = 0
    GROUP BY w.id
    HAVING count(*) FILTER (WHERE s.phase = 'Pending') <> 1 OR count(*) FILTER (WHERE s.phase = 'Finished') = 0
    LIMIT 1;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: malformed active workflow ids: %', bad_ids; END IF;

    SELECT string_agg(w.id::text, ', ' ORDER BY w.id) INTO bad_ids
    FROM project_management_workflows w
    LEFT JOIN project_management_workflows source ON source.id = w.source_template_id AND source.is_del = 0 AND source.project_id = 0
    WHERE w.is_del = 0 AND w.project_id <> 0 AND w.source_template_id <> 0 AND source.id IS NULL;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: workflows with dangling source_template_id: %', bad_ids; END IF;

    SELECT string_agg(w.id::text, ', ' ORDER BY w.id) INTO bad_ids
    FROM project_management_workflows w
    WHERE w.is_del = 0 AND w.project_id = 0 AND EXISTS (
        SELECT 1 FROM project_management_workflows other
        WHERE other.is_del = 0 AND other.project_id = 0 AND other.organization_id = w.organization_id
          AND other.name = w.name AND other.id <> w.id
    );
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: duplicate active organization workflow ids: %', bad_ids; END IF;

    SELECT string_agg(w.id::text, ', ' ORDER BY w.id) INTO bad_ids
    FROM project_management_workflows w
    WHERE w.is_del = 0 AND w.project_id = 0 AND w.is_default AND EXISTS (
        SELECT 1 FROM project_management_workflows other
        WHERE other.is_del = 0 AND other.project_id = 0 AND other.organization_id = w.organization_id
          AND other.is_default AND other.id <> w.id
    );
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: duplicate active organization default workflow ids: %', bad_ids; END IF;

    SELECT string_agg(ws.id::text, ', ' ORDER BY ws.id) INTO bad_ids
    FROM project_management_workflow_stages ws JOIN project_management_workflows w ON w.id = ws.workflow_id
    WHERE ws.is_del = 0 AND w.is_del = 0
    GROUP BY ws.workflow_id, ws.name HAVING count(*) > 1 LIMIT 1;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: duplicate active stage ids: %', bad_ids; END IF;

    SELECT string_agg(ws.id::text, ', ' ORDER BY ws.id) INTO bad_ids
    FROM project_management_workflow_stages ws JOIN project_management_workflows w ON w.id = ws.workflow_id
    WHERE ws.is_del = 0 AND w.is_del = 0 AND ws.phase = 'Pending' AND EXISTS (
        SELECT 1 FROM project_management_workflow_stages earlier
        WHERE earlier.workflow_id = ws.workflow_id AND earlier.is_del = 0
          AND (earlier.sort_order < ws.sort_order OR (earlier.sort_order = ws.sort_order AND earlier.id < ws.id))
    );
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: pending stage is not first: %', bad_ids; END IF;

    SELECT string_agg(t.id::text, ', ' ORDER BY t.id) INTO bad_ids
    FROM project_management_tasks t
    LEFT JOIN project_management_workflows w ON w.project_id = t.project_id AND w.is_del = 0
    LEFT JOIN project_management_workflow_stages ws ON ws.id = t.stage_id AND ws.is_del = 0 AND ws.workflow_id = w.id
    WHERE w.id IS NULL OR ws.id IS NULL;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: tasks with unmappable project workflow or stage ids: %', bad_ids; END IF;

    SELECT string_agg(a.id::text, ', ' ORDER BY a.id) INTO bad_ids
    FROM project_management_task_assignees a
    LEFT JOIN project_management_tasks t ON t.id = a.task_id
    LEFT JOIN project_management_workflows w ON w.project_id = t.project_id AND w.is_del = 0
    LEFT JOIN project_management_workflow_stages ws ON ws.id = a.stage_id AND ws.is_del = 0 AND ws.workflow_id = w.id
    WHERE t.id IS NULL OR w.id IS NULL OR ws.id IS NULL;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: task assignees with unmappable task, project workflow or stage ids: %', bad_ids; END IF;

    SELECT string_agg(tr.id::text, ', ' ORDER BY tr.id) INTO bad_ids
    FROM project_management_workflow_transitions tr
    LEFT JOIN project_management_workflow_stages f ON f.id = tr.from_stage_id AND f.workflow_id = tr.workflow_id AND f.is_del = 0
    LEFT JOIN project_management_workflow_stages t ON t.id = tr.to_stage_id AND t.workflow_id = tr.workflow_id AND t.is_del = 0
    WHERE tr.is_del = 0 AND (f.id IS NULL OR t.id IS NULL);
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: transitions with unmappable endpoints: %', bad_ids; END IF;

    SELECT string_agg(id::text, ', ' ORDER BY id) INTO bad_ids
    FROM project_management_project_members WHERE is_del = 0
    GROUP BY project_id, organization_user_id HAVING count(*) > 1 LIMIT 1;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 preflight: duplicate active project member ids: %', bad_ids; END IF;

INSERT INTO pm_v2_projects_staging (
    id, created_at, updated_at, is_del, organization_id, entity_id, entity_type,
    name, description, cover_url, owner_id, owner_name, priority, project_type
)
SELECT
    p.id, p.created_at, p.updated_at, p.is_del, p.organization_id, p.entity_id, p.entity_type,
    p.name, p.description, p.cover_url, p.owner_id, p.owner_name, p.priority, 'Ordinary'
FROM project_management_projects p;

-- Project workflows keep their identity only in the map; source workflows receive a template root.
INSERT INTO pm_v2_templates_staging (organization_id, name, description, template_type, status, is_default, sort_order)
SELECT w.organization_id, w.name, w.description, 'Organization', 'Active', w.is_default, 0
FROM project_management_workflows w WHERE w.is_del = 0 AND w.project_id = 0;

INSERT INTO pm_v2_workflows_staging (entity_type, entity_id, source_workflow_id, source_workflow_name, name, description, is_default)
SELECT CASE WHEN w.project_id = 0 THEN 'Template' ELSE 'Project' END,
       CASE WHEN w.project_id = 0 THEN t.id ELSE w.project_id END,
       w.id, w.name, w.name, w.description,
       CASE WHEN w.project_id = 0 THEN w.is_default ELSE true END
FROM project_management_workflows w
LEFT JOIN pm_v2_templates_staging t ON t.template_type = 'Organization' AND t.organization_id = w.organization_id AND t.name = w.name AND t.is_del = 0
WHERE w.is_del = 0;

INSERT INTO pm_v2_workflow_map (old_id, new_id, entity_type, entity_id)
SELECT w.id, n.id, n.entity_type, n.entity_id
FROM project_management_workflows w
JOIN pm_v2_workflows_staging n ON n.source_workflow_id = w.id;

INSERT INTO pm_v2_source_workflow_map (old_source_template_id, new_source_workflow_id)
SELECT old_id, new_id FROM pm_v2_workflow_map WHERE entity_type = 'Template';

UPDATE pm_v2_workflows_staging n
SET source_workflow_id = COALESCE(source_map.new_source_workflow_id, 0),
    source_workflow_name = CASE
        WHEN w.project_id <> 0 AND w.source_template_id <> 0 THEN w.source_template_name
        ELSE ''
    END
FROM pm_v2_workflow_map workflow_map
JOIN project_management_workflows w ON w.id = workflow_map.old_id
LEFT JOIN pm_v2_source_workflow_map source_map ON source_map.old_source_template_id = w.source_template_id
WHERE n.id = workflow_map.new_id;

-- v1 stages become project/template-owned stages. A v1 project is guaranteed to have one active workflow by preflight.
INSERT INTO pm_v2_stages_staging (entity_type, entity_id, name, description, phase, color, sort_order)
SELECT m.entity_type, m.entity_id, s.name, s.description, s.phase, s.color, s.sort_order
FROM project_management_workflow_stages s JOIN pm_v2_workflow_map m ON m.old_id = s.workflow_id
WHERE s.is_del = 0;
INSERT INTO pm_v2_stage_map (old_id, new_id)
SELECT s.id, n.id FROM project_management_workflow_stages s
JOIN pm_v2_workflow_map m ON m.old_id = s.workflow_id
JOIN pm_v2_stages_staging n ON n.entity_type = m.entity_type AND n.entity_id = m.entity_id AND n.name = s.name AND n.phase = s.phase
WHERE s.is_del = 0;

INSERT INTO pm_v2_workflow_stages_staging (workflow_id, stage_id)
SELECT m.new_id, sm.new_id FROM project_management_workflow_stages s
JOIN pm_v2_workflow_map m ON m.old_id = s.workflow_id JOIN pm_v2_stage_map sm ON sm.old_id = s.id
WHERE s.is_del = 0;
INSERT INTO pm_v2_workflow_transitions_staging (workflow_id, name, from_stage_id, to_stage_id, allowed_actor_types)
SELECT m.new_id, tr.name, fm.new_id, tm.new_id, tr.allowed_actor_types
FROM project_management_workflow_transitions tr JOIN pm_v2_workflow_map m ON m.old_id = tr.workflow_id
JOIN pm_v2_stage_map fm ON fm.old_id = tr.from_stage_id JOIN pm_v2_stage_map tm ON tm.old_id = tr.to_stage_id
WHERE tr.is_del = 0;

INSERT INTO pm_v2_tasks_staging (id, created_at, updated_at, is_del, entity_id, organization_id, parent_task_id, stage_id, phase, creator_id, creator_name, title, description, priority, work_hours, start_at, end_at, task_type, entity_type, workflow_id)
SELECT t.id, t.created_at, t.updated_at, t.is_del, t.project_id, t.organization_id, t.parent_task_id, sm.new_id, t.phase, t.creator_id, t.creator_name, t.title, t.description, t.priority, t.work_hours, t.start_at, t.end_at, t.task_type, 'Project', wm.new_id
FROM project_management_tasks t JOIN project_management_workflows w ON w.project_id = t.project_id AND w.is_del = 0
JOIN pm_v2_workflow_map wm ON wm.old_id = w.id JOIN pm_v2_stage_map sm ON sm.old_id = t.stage_id;
INSERT INTO pm_v2_task_assignees_staging (id, created_at, updated_at, is_del, task_id, organization_user_id, organization_user_full_name, stage_id)
SELECT a.id, a.created_at, a.updated_at, a.is_del, a.task_id, a.organization_user_id, a.organization_user_full_name, sm.new_id
FROM project_management_task_assignees a JOIN pm_v2_stage_map sm ON sm.old_id = a.stage_id;
INSERT INTO pm_v2_task_dependencies_staging SELECT * FROM project_management_task_dependencies;

    SELECT string_agg(id::text, ', ' ORDER BY id) INTO bad_ids FROM pm_v2_tasks_staging WHERE is_del = 0 AND (workflow_id IS NULL OR stage_id IS NULL OR phase IS NULL);
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 staging: unmapped active task ids: %', bad_ids; END IF;
    IF (SELECT count(*) FROM project_management_tasks) <> (SELECT count(*) FROM pm_v2_tasks_staging) THEN
        RAISE EXCEPTION 'pm v2 staging: task count mismatch';
    END IF;
    IF (SELECT count(*) FROM project_management_task_assignees) <> (SELECT count(*) FROM pm_v2_task_assignees_staging) THEN
        RAISE EXCEPTION 'pm v2 staging: task assignee count mismatch';
    END IF;
    SELECT string_agg(entity_id::text, ', ' ORDER BY entity_id) INTO bad_ids FROM pm_v2_stages_staging WHERE is_del = 0 AND entity_type = 'Project' GROUP BY entity_id HAVING count(*) FILTER (WHERE phase = 'Pending') <> 1 OR count(*) FILTER (WHERE phase = 'Finished') = 0 LIMIT 1;
    IF bad_ids IS NOT NULL THEN RAISE EXCEPTION 'pm v2 staging: project violates stage invariant: %', bad_ids; END IF;

-- Atomic cutover starts only after all staging checks have passed.
ALTER TABLE project_management_workflows RENAME TO project_management_workflows_v1;
ALTER TABLE project_management_workflow_stages RENAME TO project_management_workflow_stages_v1;
ALTER TABLE project_management_workflow_transitions RENAME TO project_management_workflow_transitions_v1;
ALTER TABLE project_management_tasks RENAME TO project_management_tasks_v1;
ALTER TABLE project_management_task_assignees RENAME TO project_management_task_assignees_v1;
ALTER TABLE project_management_task_dependencies RENAME TO project_management_task_dependencies_v1;
ALTER TABLE project_management_projects RENAME TO project_management_projects_v1;

DROP TABLE project_management_workflows_v1,
    project_management_workflow_stages_v1,
    project_management_workflow_transitions_v1,
    project_management_tasks_v1,
    project_management_task_assignees_v1,
    project_management_task_dependencies_v1,
    project_management_projects_v1;

ALTER TABLE pm_v2_projects_staging RENAME TO project_management_projects;
ALTER TABLE pm_v2_templates_staging RENAME TO project_management_templates;
ALTER TABLE pm_v2_stages_staging RENAME TO project_management_stages;
ALTER TABLE pm_v2_workflows_staging RENAME TO project_management_workflows;
ALTER TABLE pm_v2_workflow_stages_staging RENAME TO project_management_workflow_stages;
ALTER TABLE pm_v2_workflow_transitions_staging RENAME TO project_management_workflow_transitions;
ALTER TABLE pm_v2_tasks_staging RENAME TO project_management_tasks;
ALTER TABLE pm_v2_task_assignees_staging RENAME TO project_management_task_assignees;
ALTER TABLE pm_v2_task_dependencies_staging RENAME TO project_management_task_dependencies;

ALTER INDEX pm_v2_projects_staging_pkey RENAME TO project_management_projects_pkey;
ALTER INDEX pm_v2_templates_staging_pkey RENAME TO project_management_templates_pkey;
ALTER INDEX pm_v2_stages_staging_pkey RENAME TO project_management_stages_pkey;
ALTER INDEX pm_v2_workflows_staging_pkey RENAME TO project_management_workflows_pkey;
ALTER INDEX pm_v2_workflow_stages_staging_pkey RENAME TO project_management_workflow_stages_pkey;
ALTER INDEX pm_v2_workflow_transitions_staging_pkey RENAME TO project_management_workflow_transitions_pkey;
ALTER INDEX pm_v2_tasks_staging_pkey RENAME TO project_management_tasks_pkey;
ALTER INDEX pm_v2_task_assignees_staging_pkey RENAME TO project_management_task_assignees_pkey;
ALTER INDEX pm_v2_task_dependencies_staging_pkey RENAME TO project_management_task_dependencies_pkey;
ALTER SEQUENCE pm_v2_projects_staging_id_seq RENAME TO project_management_projects_id_seq;
ALTER SEQUENCE pm_v2_templates_staging_id_seq RENAME TO project_management_templates_id_seq;
ALTER SEQUENCE pm_v2_stages_staging_id_seq RENAME TO project_management_stages_id_seq;
ALTER SEQUENCE pm_v2_workflows_staging_id_seq RENAME TO project_management_workflows_id_seq;
ALTER SEQUENCE pm_v2_workflow_stages_staging_id_seq RENAME TO project_management_workflow_stages_id_seq;
ALTER SEQUENCE pm_v2_workflow_transitions_staging_id_seq RENAME TO project_management_workflow_transitions_id_seq;
ALTER SEQUENCE pm_v2_tasks_staging_id_seq RENAME TO project_management_tasks_id_seq;
ALTER SEQUENCE pm_v2_task_assignees_staging_id_seq RENAME TO project_management_task_assignees_id_seq;
ALTER SEQUENCE pm_v2_task_dependencies_staging_id_seq RENAME TO project_management_task_dependencies_id_seq;

DROP INDEX pm_v2_projects_staging_organization_id_entity_id_entity_typ_idx;
DROP INDEX pm_v2_tasks_staging_organization_id_creator_id_idx;
DROP INDEX pm_v2_tasks_staging_project_id_parent_task_id_idx;
DROP INDEX pm_v2_tasks_staging_project_id_stage_id_idx;
DROP INDEX IF EXISTS pm_v2_tasks_staging_project_id_phase_idx;
DROP INDEX pm_v2_task_assignees_staging_task_id_stage_id_organization__idx;
DROP INDEX pm_v2_task_dependencies_staging_dependency_task_id_idx;
DROP INDEX pm_v2_task_dependencies_staging_task_id_idx;
DROP INDEX idx_pm_member;

CREATE INDEX idx_pm_project_entity ON project_management_projects (organization_id, entity_id, entity_type);
CREATE INDEX idx_pm_stage_entity_name ON project_management_stages (entity_type, entity_id, name);
CREATE INDEX idx_pm_stage_entity_sort_order ON project_management_stages (entity_type, entity_id, sort_order, id);
CREATE INDEX idx_pm_stage_entity_phase ON project_management_stages (entity_type, entity_id, phase, id);
CREATE INDEX idx_pm_workflow_stage ON project_management_workflow_stages (workflow_id, stage_id);
CREATE INDEX idx_pm_transition ON project_management_workflow_transitions (workflow_id, from_stage_id, to_stage_id);
CREATE INDEX idx_pm_workflow_source ON project_management_workflows (entity_type, entity_id, source_workflow_id);
CREATE INDEX idx_pm_project_workflow_default ON project_management_workflows (entity_type, entity_id, is_default);
CREATE INDEX idx_pm_workflow_source_default ON project_management_templates (organization_id, is_default);
CREATE INDEX idx_pm_template_default_active ON project_management_templates (organization_id, template_type);
CREATE INDEX idx_pm_project_workflow_default_active ON project_management_workflows (entity_type, entity_id);
CREATE INDEX idx_pm_template_directory_name ON project_management_templates (organization_id, template_type, name);
CREATE INDEX idx_pm_template_directory_sort ON project_management_templates (organization_id, template_type, sort_order, id);
CREATE INDEX idx_pm_member ON project_management_project_members (project_id, organization_user_id);
CREATE INDEX idx_pm_task_entity_workflow_stage ON project_management_tasks (entity_type, entity_id, workflow_id, stage_id);
CREATE INDEX idx_pm_task_entity_parent ON project_management_tasks (entity_type, entity_id, parent_task_id);
CREATE INDEX idx_pm_task_creator ON project_management_tasks (organization_id, creator_id);
CREATE INDEX idx_pm_assignee ON project_management_task_assignees (task_id, stage_id, organization_user_id);
CREATE INDEX idx_pm_dependency_task ON project_management_task_dependencies (task_id);
CREATE INDEX idx_pm_dependency_related_task ON project_management_task_dependencies (dependency_task_id);



PERFORM setval(pg_get_serial_sequence('project_management_templates', 'id'), coalesce((SELECT max(id) FROM project_management_templates), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_stages', 'id'), coalesce((SELECT max(id) FROM project_management_stages), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_workflows', 'id'), coalesce((SELECT max(id) FROM project_management_workflows), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_workflow_stages', 'id'), coalesce((SELECT max(id) FROM project_management_workflow_stages), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_workflow_transitions', 'id'), coalesce((SELECT max(id) FROM project_management_workflow_transitions), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_projects', 'id'), coalesce((SELECT max(id) FROM project_management_projects), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_tasks', 'id'), coalesce((SELECT max(id) FROM project_management_tasks), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_task_assignees', 'id'), coalesce((SELECT max(id) FROM project_management_task_assignees), 1), true);
PERFORM setval(pg_get_serial_sequence('project_management_task_dependencies', 'id'), coalesce((SELECT max(id) FROM project_management_task_dependencies), 1), true);
END;
$pm_v2_migration$;
