/**
 * 任务类型 CRUD；路径与看板 `TASK_TYPES_API` 共用。
 */

import { TASK_TYPES_API } from './task';

export { TASK_TYPES_API };

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTaskTypeApi(id) {
	return `${TASK_TYPES_API}/${id}`;
}

/**
 * @param {import('@repo/shared-methods/utils').RequestFn} request
 * @returns {{
 *   list: () => Promise<unknown>,
 *   create: (body: object) => Promise<unknown>,
 *   update: (id: number|string, body: object) => Promise<unknown>,
 *   delete: (id: number|string) => Promise<unknown>,
 * }}
 */
export function createTaskTypeApi(request) {
	const hide = { isHideTips: true };

	return {
		list: () => request(TASK_TYPES_API, { method: 'GET', ...hide }),
		create: (body) => request(TASK_TYPES_API, { method: 'POST', ...hide, body }),
		update: (id, body) => request(getTaskTypeApi(id), { method: 'PUT', ...hide, body }),
		delete: (id) => request(getTaskTypeApi(id), { method: 'DELETE', ...hide }),
	};
}
