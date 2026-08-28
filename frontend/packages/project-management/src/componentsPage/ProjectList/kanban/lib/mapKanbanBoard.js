/**
 * 看板初始化映射：项目详情 workflow.stages → 列；task/page 列表 → 看板任务字段。
 * 纯映射，不改外部状态。
 */

import { PROJECT_PRIORITY } from '@repo/project-management/enum/project';
import { TASK_PRIORITY } from '@repo/project-management/enum/task';
import { mapWorkflowDetail } from '../../../Workflow/lib/mapWorkflowDetail';

const API_PRIORITY_TO_PROJECT = {
	[TASK_PRIORITY.LOW]: PROJECT_PRIORITY.LOW,
	[TASK_PRIORITY.MEDIUM]: PROJECT_PRIORITY.MEDIUM,
	[TASK_PRIORITY.HIGH]: PROJECT_PRIORITY.HIGH,
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

	return {
		...task,
		workflow_state_id: task.stage_id ?? task.workflow_state_id,
		parent_id: normalizedParentId,
		due_date: task.end_at ?? task.due_date,
		start_date: task.start_at ?? task.start_date,
		type: task.task_type?.name ?? task.type,
		priority: mapTaskPriority(task.priority),
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
