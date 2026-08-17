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

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskActivitiesApi(id) {
	return `${TASK_API}/${id}/activities`;
}

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskDependenciesApi(id) {
	return `${TASK_API}/${id}/dependencies`;
}

/**
 * @param {number|string} id
 * @param {number|string} dependencyId
 * @returns {string}
 */
export function getTaskDependencyApi(id, dependencyId) {
	return `${TASK_API}/${id}/dependencies/${dependencyId}`;
}

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskDependentsApi(id) {
	return `${TASK_API}/${id}/dependents`;
}

/**
 * @param {number|string} id
 * @param {number|string} dependencyId
 * @returns {string}
 */
export function getTaskDependentApi(id, dependencyId) {
	return `${TASK_API}/${id}/dependents/${dependencyId}`;
}

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskObserversApi(id) {
	return `${TASK_API}/${id}/observers`;
}

/**
 * @param {number|string} id
 * @param {number|string} userId
 * @returns {string}
 */
export function getTaskObserverApi(id, userId) {
	return `${TASK_API}/${id}/observers/${userId}`;
}
