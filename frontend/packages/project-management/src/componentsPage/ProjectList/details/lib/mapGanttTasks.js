/**
 * 将项目任务映射为 gantt-task-react 的 Task。
 * end 用「次日 0 点」半开区间盖满截止日；条上 name 放工期，标题放 taskTitle 供列表/hover。
 */

import { getTaskSpan, spanDays } from './taskSchedule';

/** 根条 base-primary，子条 blue-medium；今天列仍用 base-active。 */
const BAR_ROOT = 'var(--base-primary)';
const BAR_ROOT_ACTIVE = 'var(--blue-dark)';
const BAR_CHILD = 'var(--blue-medium)';
const BAR_CHILD_ACTIVE = 'var(--blue-dark)';

/**
 * @param {string} color
 * @param {string} active
 * @returns {{ backgroundColor: string, backgroundSelectedColor: string, progressColor: string, progressSelectedColor: string }}
 */
function barStyles(color, active) {
	return {
		backgroundColor: color,
		backgroundSelectedColor: active,
		progressColor: color,
		progressSelectedColor: active,
	};
}

/**
 * 树序展开；有子任务的节点 type=project，便于左侧折叠。
 * @param {Array} tasks
 * @param {{ collapsedIds?: Set<string>, formatDuration: (days: number) => string }} options
 * @returns {Array}
 */
export function mapGanttTasks(tasks, { collapsedIds = new Set(), formatDuration }) {
	if (!tasks?.length) return [];

	const byId = new Map(tasks.map((task) => [task.id, task]));
	const childrenOf = new Map();
	tasks.forEach((task) => {
		const parentId = task.parent_id;
		if (parentId == null || !byId.has(parentId) || parentId === task.id) return;
		const list = childrenOf.get(parentId) || [];
		list.push(task);
		childrenOf.set(parentId, list);
	});

	const roots = tasks.filter((task) => {
		const parentId = task.parent_id;
		return parentId == null || !byId.has(parentId) || parentId === task.id;
	});

	const sortByStart = (list) =>
		[...list].sort((a, b) => getTaskSpan(a).start.valueOf() - getTaskSpan(b).start.valueOf());

	/** @type {Array} */
	const result = [];
	let order = 1;

	/**
	 * @param {object} task
	 * @param {string|undefined} projectId
	 */
	const walk = (task, projectId) => {
		const span = getTaskSpan(task);
		const childList = childrenOf.get(task.id) || [];
		const hasChildren = childList.length > 0;
		const id = String(task.id);
		const days = spanDays(task);
		const title = task.title || `#${task.id}`;
		const color = projectId ? BAR_CHILD : BAR_ROOT;
		const active = projectId ? BAR_CHILD_ACTIVE : BAR_ROOT_ACTIVE;

		result.push({
			id,
			/** 插件条形文案读 name；工期上条，标题只在左侧/hover */
			name: formatDuration(days),
			taskTitle: title,
			type: hasChildren ? 'project' : 'task',
			start: span.start.toDate(),
			end: span.end.add(1, 'day').startOf('day').toDate(),
			progress: 0,
			isDisabled: true,
			project: projectId,
			hideChildren: hasChildren ? collapsedIds.has(id) : undefined,
			displayOrder: order++,
			durationDays: days,
			displayStart: span.start.toDate(),
			displayEnd: span.end.toDate(),
			styles: barStyles(color, active),
		});

		sortByStart(childList).forEach((child) => walk(child, id));
	};

	sortByStart(roots).forEach((task) => walk(task, undefined));
	return result;
}
