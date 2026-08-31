/**
 * 任务接口：我的任务分页 POST `/project_management/task/my`；
 * 看板任务分页 POST `/project_management/task/page`（filter 含 project_id / stage_id）；
 * 看板新建 POST `/project_management/task`；更新 / 流转仍走 `/tasks`（待迁移）。
 */

export const TASK_API = '/tasks';
export const TASK_TYPES_API = '/task-types';

/** 我的任务分页：POST，body `{ page, limit, filter }`，data 为 `{ list, total }` */
export const TASK_MY_API = '/project_management/task/my';

/** 项目看板任务分页：POST，body `{ page, limit, filter }`，filter 含 project_id、stage_id 等 */
export const TASK_PAGE_API = '/project_management/task/page';

/** 看板新建任务：POST，body 含 project_id、task_type_id、priority（字符串）等 */
export const TASK_CREATE_API = '/project_management/task';

/**
 * @param {import('@repo/shared-methods/utils').RequestFn} request
 * @returns {{
 *   page: (body: object) => Promise<{ list?: object[], total?: number }|unknown>,
 *   pageByProject: (body: object) => Promise<{ list?: object[], total?: number }|unknown>,
 *   create: (body: object) => Promise<unknown>,
 * }}
 */
export function createTaskApi(request) {
	const hide = { isHideTips: true };
	return {
		page: (body) => request(TASK_MY_API, { method: 'POST', ...hide, body }),
		pageByProject: (body) => request(TASK_PAGE_API, { method: 'POST', ...hide, body }),
		create: (body) => request(TASK_CREATE_API, { method: 'POST', ...hide, body }),
	};
}

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
