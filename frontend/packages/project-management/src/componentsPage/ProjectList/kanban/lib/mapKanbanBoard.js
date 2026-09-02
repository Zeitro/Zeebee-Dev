/**
 * 看板初始化映射：项目详情 workflow.stages → 列；task/page 列表 → 看板任务字段。
 * 纯映射，不改外部状态。
 */

import { PROJECT_PRIORITY } from '@repo/project-management/enum/project';
import { TASK_PRIORITY, toTaskPriorityValue } from '@repo/project-management/enum/task';
import { mapWorkflowDetail } from '../../../Workflow/lib/mapWorkflowDetail';
import { getAssigneeIds, toApiDate } from '../components/helpers';

const API_PRIORITY_TO_PROJECT = {
	[TASK_PRIORITY.LOW]: PROJECT_PRIORITY.LOW,
	[TASK_PRIORITY.MEDIUM]: PROJECT_PRIORITY.MEDIUM,
	[TASK_PRIORITY.HIGH]: PROJECT_PRIORITY.HIGH,
	[TASK_PRIORITY.URGENT]: PROJECT_PRIORITY.URGENT,
	[TASK_PRIORITY.CRITICAL]: PROJECT_PRIORITY.URGENT,
};

/**
 * 接口 priority（字符串）→ 看板卡片用的数字优先级。
 * @param {string|number|null|undefined} priority
 * @returns {number}
 */
function mapTaskPriority(priority) {
	if (typeof priority === 'number' && !Number.isNaN(priority)) return priority;
	if (typeof priority === 'string' && priority) {
		return API_PRIORITY_TO_PROJECT[priority] ?? PROJECT_PRIORITY.MEDIUM;
	}
	return PROJECT_PRIORITY.MEDIUM;
}

/**
 * POST `/project_management/task/page` 单条 → 看板内部任务结构。
 * @param {object|null|undefined} task
 * @returns {object|null}
 */
export function mapKanbanTask(task) {
	if (!task || typeof task !== 'object') return null;

	const parentId = task.parent_task_id ?? task.parent_id;
	const normalizedParentId =
		parentId === null || parentId === undefined || parentId === 0 ? null : parentId;

	const assigneeIds = Array.isArray(task.assignees)
		? task.assignees
				.map((item) => item.organization_user_id ?? item.user_id)
				.filter((id) => id != null)
		: Array.isArray(task.assignee_ids)
			? task.assignee_ids
			: task.assignee_id
				? [task.assignee_id]
				: [];

	return {
		...task,
		workflow_state_id: task.stage_id ?? task.workflow_state_id,
		parent_id: normalizedParentId,
		due_date: task.end_at ?? task.due_date,
		start_date: task.start_at ?? task.start_date,
		type: task.task_type?.name ?? task.type,
		task_type_id: task.task_type?.id ?? task.task_type_id,
		priority: mapTaskPriority(task.priority),
		assignee_ids: assigneeIds,
		assignee_id: assigneeIds[0] ?? task.assignee_id ?? null,
		estimated_hours: task.work_hours ?? task.estimated_hours,
	};
}

/**
 * GET `/project_management/task/:id` → 看板详情任务结构（含 assignees / attachments / dependencies / children）。
 * @param {object|null|undefined} task
 * @returns {object|null}
 */
export function mapKanbanTaskDetail(task) {
	const mapped = mapKanbanTask(task);
	if (!mapped) return null;

	const children = Array.isArray(task.children)
		? task.children.map((item) => mapKanbanTask(item)).filter(Boolean)
		: mapped.children;

	return {
		...mapped,
		attachments: task.attachments ?? mapped.attachments,
		dependencies: task.dependencies ?? mapped.dependencies,
		children,
	};
}

/**
 * 单条依赖 → PUT body 项；新建项 `id` 传 0。
 * @param {object|null|undefined} dep
 * @returns {object|null}
 */
export function mapDependencyForUpdate(dep) {
	if (!dep || dep.dependency_task_id == null || !dep.dependency_type) return null;
	return {
		id: Number(dep.id) || 0,
		task_id: Number(dep.task_id),
		dependency_task_id: Number(dep.dependency_task_id),
		dependency_type: dep.dependency_type,
		dependency_rule: dep.dependency_rule,
	};
}

/**
 * 单条附件 → PUT body 项；新建项 `file_id` 传 0。
 * @param {object|null|undefined} attachment
 * @returns {object|null}
 */
export function mapAttachmentForUpdate(attachment) {
	if (!attachment?.file_url) return null;
	return {
		file_id: Number(attachment.file_id) || 0,
		file_name: attachment.file_name,
		file_url: attachment.file_url,
		storage_type: attachment.storage_type || 'Aws S3',
	};
}

/**
 * 看板任务详情 → PUT `/project_management/task/:id` body。
 * @param {object} task 当前详情
 * @param {object} [overrides] 表单或局部变更
 * @param {{ taskTypes?: Array }} [options]
 * @returns {object}
 */
export function buildTaskUpdateBody(task, overrides = {}, options = {}) {
	const { taskTypes = [] } = options;

	let taskTypeId = overrides.task_type_id ?? task.task_type_id ?? task.task_type?.id;
	if (overrides.type != null) {
		const matched = taskTypes.find((item) => item.name === overrides.type);
		if (matched?.id) taskTypeId = matched.id;
	}

	const assigneeIds =
		overrides.assignee_organization_user_ids ??
		(overrides.assignee_ids != null
			? overrides.assignee_ids.map((id) => Number(id)).filter(Boolean)
			: getAssigneeIds(task));

	const dependencies = (overrides.dependencies ?? task.dependencies ?? [])
		.map(mapDependencyForUpdate)
		.filter(Boolean);

	const attachments = (overrides.attachments ?? task.attachments ?? [])
		.map(mapAttachmentForUpdate)
		.filter(Boolean);

	const workHours =
		overrides.work_hours ??
		overrides.estimated_hours ??
		task.work_hours ??
		task.estimated_hours;

	return {
		title: overrides.title ?? task.title,
		description: overrides.description ?? task.description ?? '',
		priority: toTaskPriorityValue(overrides.priority ?? task.priority),
		task_type_id: taskTypeId ? Number(taskTypeId) : undefined,
		start_at: toApiDate(
			overrides.start_at ?? overrides.start_date ?? task.start_date ?? task.start_at,
		),
		end_at: toApiDate(overrides.end_at ?? overrides.due_date ?? task.due_date ?? task.end_at),
		work_hours: workHours != null && workHours !== '' ? Number(workHours) : undefined,
		assignee_organization_user_ids: assigneeIds,
		dependencies,
		attachments,
	};
}

/**
 * GET `/project_management/project/:id` → 看板 project / workflow。
 * @param {object|null|undefined} data
 * @returns {{ project: object|null, workflow: object|null, taskTotal: number }}
 */
export function mapKanbanProjectDetail(data) {
	if (!data || typeof data !== 'object') {
		return { project: null, workflow: null, taskTotal: 0 };
	}

	const workflow = data.workflow ? mapWorkflowDetail(data.workflow) : null;
	const taskTotal = Number(data.task_count ?? data.total_task_count ?? data.tasks_count) || 0;

	return {
		project: data,
		workflow,
		taskTotal,
	};
}
