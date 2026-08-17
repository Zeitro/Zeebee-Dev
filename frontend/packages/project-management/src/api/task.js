/**
 * 任务接口：列表 GET query；创建 POST；更新 PUT；拖拽改状态走 transition。
 */

export const TASK_API = '/tasks';
export const TASK_TYPES_API = '/task-types';

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskApi(id) {
	return `${TASK_API}/${id}`;
}

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskTransitionApi(id) {
	return `${TASK_API}/${id}/transition`;
}
