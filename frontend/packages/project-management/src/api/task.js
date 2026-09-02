/**
 * 任务接口：我的任务分页 POST `/project_management/task/my`；
 * 看板任务分页 POST `/project_management/task/page`（filter 含 project_id / stage_id）；
 * 看板新建 POST `/project_management/task`；详情 GET `/project_management/task/:id?project_id=`；
 * 依赖候选 GET `/project_management/task/candidate?project_id=&task_id=&keyword=`；
 * 编辑 PUT `/project_management/task/:id`（含 dependencies 等全量字段）；
 * 删除 DELETE `/project_management/task/:id?project_id=`；
 * 流转 POST `/project_management/task/:id/transition`，body `{ project_id, target_stage_id }`。
 */

export const TASK_API = '/tasks';
export const TASK_TYPES_API = '/task-types';

/** 我的任务分页：POST，body `{ page, limit, filter }`，data 为 `{ list, total }` */
export const TASK_MY_API = '/project_management/task/my';

/** 项目看板任务分页：POST，body `{ page, limit, filter }`，filter 含 project_id、stage_id 等 */
export const TASK_PAGE_API = '/project_management/task/page';

/** 看板新建任务：POST，body 含 project_id、task_type_id、priority（字符串）等 */
export const TASK_CREATE_API = '/project_management/task';

/** 依赖候选任务：GET，query `{ project_id, task_id, keyword? }`，data 为 `{ list }`（`id` / `title` / `phase`） */
export const TASK_CANDIDATE_API = `${TASK_CREATE_API}/candidate`;

/**
 * 任务详情 query：`project_id` 必填。
 * @param {number|string|undefined|null} projectId
 * @returns {string}
 */
export function buildTaskDetailQuery(projectId) {
	const search = new URLSearchParams();
	if (projectId !== undefined && projectId !== null && projectId !== '') {
		search.set('project_id', String(projectId));
	}
	const qs = search.toString();
	return qs ? `?${qs}` : '';
}

/**
 * 依赖候选任务 query：`project_id`、`task_id` 必填，`keyword` 可选。
 * @param {{ project_id?: number|string, task_id?: number|string, keyword?: string }} params
 * @returns {string}
 */
export function buildTaskCandidateQuery(params = {}) {
	const search = new URLSearchParams();
	Object.entries(params).forEach(([key, value]) => {
		if (value === undefined || value === null || value === '') return;
		search.set(key, String(value));
	});
	const qs = search.toString();
	return qs ? `?${qs}` : '';
}

/**
 * 任务详情：GET `/project_management/task/:id?project_id=`。
 * @param {number|string} id
 * @param {number|string|undefined|null} projectId
 * @returns {string}
 */
export function getTaskDetailApi(id, projectId) {
	return `${TASK_CREATE_API}/${id}${buildTaskDetailQuery(projectId)}`;
}

/**
 * 任务删除：DELETE `/project_management/task/:id?project_id=`。
 * @param {number|string} id
 * @param {number|string|undefined|null} projectId
 * @returns {string}
 */
export function getTaskDeleteApi(id, projectId) {
	return getTaskDetailApi(id, projectId);
}

/**
 * @param {import('@repo/shared-methods/utils').RequestFn} request
 * @returns {{
 *   page: (body: object) => Promise<{ list?: object[], total?: number }|unknown>,
 *   pageByProject: (body: object) => Promise<{ list?: object[], total?: number }|unknown>,
 *   create: (body: object) => Promise<unknown>,
 *   detail: (id: number|string, projectId: number|string) => Promise<unknown>,
 *   candidate: (params: { project_id: number|string, task_id: number|string, keyword?: string }) => Promise<{ list?: object[] }|unknown>,
 *   update: (id: number|string, body: object) => Promise<unknown>,
 *   delete: (id: number|string, projectId: number|string) => Promise<unknown>,
 *   transition: (id: number|string, body: { project_id: number|string, target_stage_id: number|string }) => Promise<unknown>,
 * }}
 */
export function createTaskApi(request) {
	const hide = { isHideTips: true };
	return {
		page: (body) => request(TASK_MY_API, { method: 'POST', ...hide, body }),
		pageByProject: (body) => request(TASK_PAGE_API, { method: 'POST', ...hide, body }),
		create: (body) => request(TASK_CREATE_API, { method: 'POST', ...hide, body }),
		detail: (id, projectId) =>
			request(getTaskDetailApi(id, projectId), { method: 'GET', ...hide }),
		candidate: (params) =>
			request(`${TASK_CANDIDATE_API}${buildTaskCandidateQuery(params)}`, {
				method: 'GET',
				...hide,
			}),
		update: (id, body) => request(`${TASK_CREATE_API}/${id}`, { method: 'PUT', ...hide, body }),
		delete: (id, projectId) =>
			request(getTaskDeleteApi(id, projectId), { method: 'DELETE', ...hide }),
		transition: (id, body) =>
			request(getTaskTransitionApi(id), { method: 'POST', ...hide, body }),
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
	return `${TASK_CREATE_API}/${id}/transition`;
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
