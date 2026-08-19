/**
 * 工作流接口路径与 `createWorkflowApi` 工厂；列表页 / 编辑器共用。
 */

export { createRequest } from '@repo/shared-methods/utils';

import { PROJECT_API } from './project';

export const WORKFLOW_API = '/workflows';

/**
 * @param {Record<string, string|number|undefined|null>} params
 * @returns {string}
 */
export function buildWorkflowQuery(params) {
	const search = new URLSearchParams();
	Object.entries(params || {}).forEach(([key, value]) => {
		if (value === undefined || value === null || value === '') return;
		search.set(key, String(value));
	});
	const qs = search.toString();
	return qs ? `?${qs}` : '';
}

/**
 * @param {number|string} projectId
 * @returns {string}
 */
export function getProjectWorkflowApi(projectId) {
	return `${PROJECT_API}/${projectId}/workflow`;
}

/**
 * @param {import('@repo/shared-methods/utils').RequestFn} request
 * @returns {{
 *   list: () => Promise<unknown>,
 *   get: (id: number|string) => Promise<unknown>,
 *   create: (body: object) => Promise<unknown>,
 *   update: (id: number|string, body: object) => Promise<unknown>,
 *   delete: (id: number|string, userId: number|string) => Promise<unknown>,
 *   getForProject: (projectId: number|string) => Promise<unknown>,
 *   addState: (workflowId: number|string, body: object) => Promise<unknown>,
 *   updateState: (workflowId: number|string, stateId: number|string, body: object) => Promise<unknown>,
 *   deleteState: (workflowId: number|string, stateId: number|string, userId: number|string) => Promise<unknown>,
 *   reorderStates: (workflowId: number|string, stateIds: Array<number|string>, userId: number|string) => Promise<unknown>,
 *   addTransition: (workflowId: number|string, body: object) => Promise<unknown>,
 *   updateTransition: (workflowId: number|string, transitionId: number|string, body: object) => Promise<unknown>,
 *   deleteTransition: (workflowId: number|string, transitionId: number|string, userId: number|string) => Promise<unknown>,
 * }}
 */
export function createWorkflowApi(request) {
	const hide = { isHideTips: true };

	return {
		list: () => request(WORKFLOW_API, { method: 'GET', ...hide }),
		get: (id) => request(`${WORKFLOW_API}/${id}`, { method: 'GET', ...hide }),
		create: (body) => request(WORKFLOW_API, { method: 'POST', ...hide, body }),
		update: (id, body) => request(`${WORKFLOW_API}/${id}`, { method: 'PUT', ...hide, body }),
		delete: (id, userId) =>
			request(`${WORKFLOW_API}/${id}${buildWorkflowQuery({ user_id: userId })}`, {
				method: 'DELETE',
				...hide,
			}),
		getForProject: (projectId) =>
			request(getProjectWorkflowApi(projectId), { method: 'GET', ...hide }),
		addState: (workflowId, body) =>
			request(`${WORKFLOW_API}/${workflowId}/states`, { method: 'POST', ...hide, body }),
		updateState: (workflowId, stateId, body) =>
			request(`${WORKFLOW_API}/${workflowId}/states/${stateId}`, {
				method: 'PUT',
				...hide,
				body,
			}),
		deleteState: (workflowId, stateId, userId) =>
			request(
				`${WORKFLOW_API}/${workflowId}/states/${stateId}${buildWorkflowQuery({ user_id: userId })}`,
				{ method: 'DELETE', ...hide },
			),
		reorderStates: (workflowId, stateIds, userId) =>
			request(`${WORKFLOW_API}/${workflowId}/states/reorder`, {
				method: 'PUT',
				...hide,
				body: { state_ids: stateIds, user_id: userId },
			}),
		addTransition: (workflowId, body) =>
			request(`${WORKFLOW_API}/${workflowId}/transitions`, { method: 'POST', ...hide, body }),
		updateTransition: (workflowId, transitionId, body) =>
			request(`${WORKFLOW_API}/${workflowId}/transitions/${transitionId}`, {
				method: 'PUT',
				...hide,
				body,
			}),
		deleteTransition: (workflowId, transitionId, userId) =>
			request(
				`${WORKFLOW_API}/${workflowId}/transitions/${transitionId}${buildWorkflowQuery({ user_id: userId })}`,
				{ method: 'DELETE', ...hide },
			),
	};
}
