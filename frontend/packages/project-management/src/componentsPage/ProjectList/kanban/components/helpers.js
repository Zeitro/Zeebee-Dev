/**
 * 看板纯函数：列表解包、query、执行者 id、成员展示；不改外部状态。
 */

import { dayjs } from '@repo/shared-methods/utils';
import { PROJECT_PRIORITY } from '@repo/project-management/enum';

/**
 * 兼容 `{ list }` / `{ data }` / 裸数组。
 * @param {*} res
 * @returns {Array}
 */
export function unwrapList(res) {
	if (Array.isArray(res)) return res;
	if (Array.isArray(res?.list)) return res.list;
	if (Array.isArray(res?.data)) return res.data;
	return [];
}

/**
 * @param {Record<string, string|number|undefined|null>} params
 * @returns {string}
 */
export function buildQuery(params) {
	const search = new URLSearchParams();
	Object.entries(params).forEach(([key, value]) => {
		if (value === undefined || value === null || value === '') return;
		search.set(key, String(value));
	});
	const qs = search.toString();
	return qs ? `?${qs}` : '';
}

/**
 * 优先 `assignee_ids`，否则回落 `assignee_id`。
 * @param {object|null|undefined} task
 * @returns {number[]}
 */
export function getAssigneeIds(task) {
	if (Array.isArray(task?.assignee_ids) && task.assignee_ids.length) {
		return task.assignee_ids;
	}
	if (task?.assignee_id) return [task.assignee_id];
	return [];
}

/**
 * @param {object|null|undefined} member
 * @param {number|string} [fallbackId]
 * @returns {string}
 */
export function getMemberLabel(member, fallbackId) {
	return member?.name || member?.email || (fallbackId != null ? String(fallbackId) : '');
}

/**
 * @param {object|null|undefined} member
 * @param {number|string} [fallbackId]
 * @returns {string}
 */
export function getMemberInitial(member, fallbackId) {
	const label = getMemberLabel(member, fallbackId);
	return label ? label.charAt(0).toUpperCase() : '?';
}

/**
 * DatePicker / dayjs 转后端 ISO；空值不传。
 * @param {*} value
 * @returns {string|undefined}
 */
export function toIsoDateTime(value) {
	if (!value) return undefined;
	const parsed = typeof value.format === 'function' ? value : dayjs(value);
	if (!parsed || !parsed.isValid?.()) return undefined;
	return parsed.format('YYYY-MM-DDTHH:mm:ssZ');
}

/**
 * 卡片不展示「中」优先级（与 demo 一致：priority === 1 隐藏）。
 * @param {number|string|null|undefined} priority
 * @returns {boolean}
 */
export function shouldShowCardPriority(priority) {
	if (priority === null || priority === undefined || priority === '') return false;
	return Number(priority) !== PROJECT_PRIORITY.MEDIUM;
}

/**
 * @param {Array} tasks
 * @param {number|string} taskId
 * @returns {number}
 */
export function countSubtasks(tasks, taskId) {
	if (!Array.isArray(tasks) || taskId == null) return 0;
	return tasks.filter((item) => item?.parent_id === taskId).length;
}

/**
 * @param {Array} tasks
 * @param {number|string} parentId
 * @returns {Array}
 */
export function getSubtasks(tasks, parentId) {
	if (!Array.isArray(tasks) || parentId == null) return [];
	return tasks.filter((item) => item?.parent_id === parentId);
}
