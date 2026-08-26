/**
 * 任务类型 CRUD。
 * 全量列表：GET `/project_management/task_type/all`，无分页，data 为 `{ list }`。
 * 新建：POST `/project_management/task_type`，body `{ name, description, color, sort_order }`。
 * 更新：PUT `/project_management/task_type/:id`，body 同上。
 * 删除：DELETE `/project_management/task_type/:id`。
 */

import { TASK_TYPES_API } from './task';

export { TASK_TYPES_API };

export const TASK_TYPE_API = '/project_management/task_type';

/** 全量任务类型：GET，响应 data 为 `{ list }` */
export const TASK_TYPE_ALL_API = `${TASK_TYPE_API}/all`;

/**
 * 详情 / 更新 / 删除：`/project_management/task_type/:id`。
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskTypeApi(id) {
	return `${TASK_TYPE_API}/${id}`;
}

/**
 * @param {import('@repo/shared-methods/utils').RequestFn} request
 * @returns {{
 *   list: () => Promise<{ list?: object[] }|unknown>,
 *   create: (body: object) => Promise<unknown>,
 *   update: (id: number|string, body: object) => Promise<unknown>,
 *   delete: (id: number|string) => Promise<unknown>,
 * }}
 */
export function createTaskTypeApi(request) {
	const hide = { isHideTips: true };

	return {
		list: () => request(TASK_TYPE_ALL_API, { method: 'GET', ...hide }),
		create: (body) => request(TASK_TYPE_API, { method: 'POST', ...hide, body }),
		update: (id, body) => request(getTaskTypeApi(id), { method: 'PUT', ...hide, body }),
		delete: (id) => request(getTaskTypeApi(id), { method: 'DELETE', ...hide }),
	};
}
