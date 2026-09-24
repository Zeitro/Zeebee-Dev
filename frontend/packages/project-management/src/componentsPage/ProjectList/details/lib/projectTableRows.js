/**
 * 任务表格行：关键字过滤后拼成树。排序由列 `sorter` 处理，这里不改顺序。
 */
import { parseTaskDate } from './taskSchedule';

/**
 * 列排序：空值视为 0，升降序由 Table 自己翻转。
 * @param {string} field
 * @returns {(left: object, right: object) => number}
 */
export function taskColumnSorter(field) {
	return (left, right) => {
		const leftValue = columnSortValue(left, field);
		const rightValue = columnSortValue(right, field);
		return leftValue - rightValue;
	};
}

/**
 * @param {object} task
 * @param {string} field
 * @returns {number}
 */
function columnSortValue(task, field) {
	if (field === 'priority') {
		const value = Number(task?.priority);
		return Number.isNaN(value) ? 0 : value;
	}
	return parseTaskDate(task?.[field])?.valueOf() || 0;
}

/**
 * @param {Array} tasks
 * @param {string} keyword
 * @returns {Array}
 */
export function filterProjectTasks(tasks, keyword) {
	const query = keyword.trim().toLowerCase();
	if (!query) return tasks;
	return tasks.filter((task) =>
		String(task.title || '')
			.toLowerCase()
			.includes(query),
	);
}

/**
 * @param {Array} tasks
 * @returns {Array}
 */
export function buildProjectTaskTree(tasks) {
	const nodes = new Map(tasks.map((task) => [task.id, { ...task, children: [] }]));
	const roots = [];
	nodes.forEach((node) => {
		const parent = node.parent_id ? nodes.get(node.parent_id) : null;
		if (parent && parent.id !== node.id) parent.children.push(node);
		else roots.push(node);
	});
	const strip = (list) =>
		list.map((node) =>
			node.children.length
				? { ...node, children: strip(node.children) }
				: { ...node, children: undefined },
		);
	return strip(roots);
}

/**
 * @param {Array} tasks
 * @returns {Map<number, number>}
 */
export function countChildTasks(tasks) {
	const counts = new Map();
	tasks.forEach((task) => {
		if (task?.parent_id == null) return;
		counts.set(task.parent_id, (counts.get(task.parent_id) || 0) + 1);
	});
	return counts;
}
