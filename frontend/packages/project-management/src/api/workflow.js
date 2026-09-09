/**
 * 工作流接口路径与 `createWorkflowApi` 工厂；列表页 / 编辑器共用。
 */

export { createRequest } from '@repo/shared-methods/utils';

export const WORKFLOW_API = '/workflows';

/** 新建工作流：POST，body `{ name, description }`，响应 data 为工作流对象（含 stages / transitions） */
export const WORKFLOW_PM_API = '/project_management/workflow';

/** 全量工作流（列表页 / 创建项目下拉）：GET，响应 data 为 `{ list }` */
export const WORKFLOW_ALL_API = `${WORKFLOW_PM_API}/all`;

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
 * 项目维度工作流：GET / PUT `/project_management/project/:id/workflow`。
 * @param {number|string} projectId
 * @returns {string}
 */
export function getProjectWorkflowApi(projectId) {
	return `/project_management/project/${projectId}/workflow`;
}

/**
 * @param {number|string|null|undefined} projectId
 * @returns {boolean}
 */
export function isValidProjectId(projectId) {
	return (
		projectId != null &&
		projectId !== '' &&
		!Number.isNaN(Number(projectId)) &&
		Number(projectId) !== 0
	);
}

/**
 * 调用方显式 projectId 优先，否则用工厂上下文 projectId。
 * @param {number|string|null|undefined} explicitProjectId
 * @param {number|string|null|undefined} contextProjectId
 * @returns {number|string|null}
 */
function resolveProjectId(explicitProjectId, contextProjectId) {
	const candidate = explicitProjectId ?? contextProjectId;
	return isValidProjectId(candidate) ? candidate : null;
}

/**
 * 详情 GET / 全量更新 PUT / 删除：`/project_management/workflow/:id`。
 * GET 的 data 含 `stages`（每项可嵌套 `transitions`）。
 * PUT body：`{ name, description, project_id, stages }`（添加 / 编辑 / 删除状态与流转均改 stages 后整份提交）。
 * @param {number|string} id
 * @returns {string}
 */
export function getWorkflowApi(id) {
	return `${WORKFLOW_PM_API}/${id}`;
}

/**
 * 设为默认：PUT `/project_management/workflow/:id/default`，无 body，data 为空对象。
 * @param {number|string} id
 * @returns {string}
 */
export function getWorkflowDefaultApi(id) {
	return `${getWorkflowApi(id)}/default`;
}

/**
 * 交换两个阶段顺序：PUT `/project_management/workflow/:id/swap_stage_order`。
 * body：`{ from_stage_id, to_stage_id }`。
 * @param {number|string} id
 * @returns {string}
 */
export function getWorkflowSwapStageOrderApi(id) {
	return `${getWorkflowApi(id)}/swap_stage_order`;
}

/**
 * @param {import('@repo/shared-methods/utils').RequestFn} request
 * @param {{ projectId?: number|string }} [options] 项目上下文；有值时 get / update 走 `/project_management/project/:id/workflow`
 * @returns {{
 *   list: () => Promise<{ list?: object[] }|unknown>,
 *   get: (id: number|string, projectId?: number|string) => Promise<object|unknown>,
 *   create: (body: { name: string, description?: string }) => Promise<object|unknown>,
 *   update: (id: number|string, body: { name: string, description?: string, project_id?: number, stages: object[] }, projectId?: number|string) => Promise<unknown>,
 *   setDefault: (id: number|string) => Promise<unknown>,
 *   delete: (id: number|string) => Promise<unknown>,
 *   getForProject: (projectId: number|string) => Promise<unknown>,
 *   addState: (workflowId: number|string, body: object) => Promise<unknown>,
 *   updateState: (workflowId: number|string, stateId: number|string, body: object) => Promise<unknown>,
 *   deleteState: (workflowId: number|string, stateId: number|string, userId: number|string) => Promise<unknown>,
 *   reorderStates: (workflowId: number|string, stateIds: Array<number|string>, userId: number|string) => Promise<unknown>,
 *   swapStageOrder: (workflowId: number|string, body: { from_stage_id: number|string, to_stage_id: number|string }) => Promise<unknown>,
 *   addTransition: (workflowId: number|string, body: object) => Promise<unknown>,
 *   updateTransition: (workflowId: number|string, transitionId: number|string, body: object) => Promise<unknown>,
 *   deleteTransition: (workflowId: number|string, transitionId: number|string, userId: number|string) => Promise<unknown>,
 * }}
 */
export function createWorkflowApi(request, options = {}) {
	const hide = { isHideTips: true };
	const contextProjectId = options.projectId;

	return {
		list: () => request(WORKFLOW_ALL_API, { method: 'GET', ...hide }),
		get: (id, projectId) => {
			const pid = resolveProjectId(projectId, contextProjectId);
			const url = pid ? getProjectWorkflowApi(pid) : getWorkflowApi(id);
			return request(url, { method: 'GET', ...hide });
		},
		create: (body) => request(WORKFLOW_PM_API, { method: 'POST', ...hide, body }),
		update: (id, body, projectId) => {
			const pid = resolveProjectId(projectId, contextProjectId);
			const url = pid ? getProjectWorkflowApi(pid) : getWorkflowApi(id);
			return request(url, { method: 'PUT', ...hide, body });
		},
		setDefault: (id) => request(getWorkflowDefaultApi(id), { method: 'PUT', ...hide }),
		delete: (id) => request(getWorkflowApi(id), { method: 'DELETE', ...hide }),
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
		swapStageOrder: (workflowId, body) =>
			request(getWorkflowSwapStageOrderApi(workflowId), { method: 'PUT', ...hide, body }),
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
