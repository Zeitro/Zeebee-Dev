/**
 * 项目成员分页：POST `/project_management/project_member/page`，
 * body `{ page, limit, filter: { project_id, keyword, role } }`，
 * 响应 data 为 `{ list, total, page, limit, total_page }`。
 * 邀请候选：GET `/project_management/project_member/candidate`，
 * query `{ project_id, keyword, purpose }`，响应 data 为 `{ list }`（`organization_user_id` / `full_name` / `email`）。
 * 批量添加：POST `/project_management/project_member`，
 * body `{ project_id, members: [{ organization_user_id, role }] }`。
 * 改角色：PUT `/project_management/project_member/:id/role`，
 * body `{ project_id, role }`。
 * 删除：DELETE `/project_management/project_member/:id`，
 * body `{ project_id }`。
 */

export const PROJECT_MEMBER_API = '/project_management/project_member';
export const PROJECT_MEMBER_PAGE_API = `${PROJECT_MEMBER_API}/page`;
export const PROJECT_MEMBER_CANDIDATE_API = `${PROJECT_MEMBER_API}/candidate`;

/** 邀请成员弹窗：拉取可添加的租户成员 */
export const PROJECT_MEMBER_CANDIDATE_PURPOSE_ADD_MEMBER = 'AddMember';

/** 邀请成员默认角色 */
export const PROJECT_MEMBER_DEFAULT_ROLE = 'TaskExecutor';

/** 看板成员列表一次拉全量上限 */
export const PROJECT_MEMBER_LIST_PAGE_SIZE = 200;

/**
 * @param {{ project_id?: number|string, keyword?: string, purpose?: string }} params
 * @returns {string}
 */
export function buildProjectMemberCandidateQuery(params = {}) {
	const search = new URLSearchParams();
	Object.entries(params).forEach(([key, value]) => {
		if (value === undefined || value === null || value === '') return;
		search.set(key, String(value));
	});
	const qs = search.toString();
	return qs ? `?${qs}` : '';
}

/**
 * 改成员角色：PUT `/project_management/project_member/:id/role`。
 * @param {number|string} id 项目成员记录 id
 * @returns {string}
 */
export function getProjectMemberRoleApi(id) {
	return `${PROJECT_MEMBER_API}/${id}/role`;
}

/**
 * 删除成员：DELETE `/project_management/project_member/:id`。
 * @param {number|string} id 项目成员记录 id
 * @returns {string}
 */
export function getProjectMemberByIdApi(id) {
	return `${PROJECT_MEMBER_API}/${id}`;
}

/**
 * @param {import('@repo/shared-methods/utils').RequestFn} request
 * @returns {{
 *   page: (body: object) => Promise<{ list?: object[], total?: number }|unknown>,
 *   candidate: (params: { project_id: number|string, keyword?: string, purpose?: string }) => Promise<{ list?: object[] }|unknown>,
 *   addMembers: (body: { project_id: number|string, members: Array<{ organization_user_id: number, role: string }> }) => Promise<unknown>,
 *   updateRole: (id: number|string, body: { project_id: number|string, role: string }) => Promise<unknown>,
 *   delete: (id: number|string, body: { project_id: number|string }) => Promise<unknown>,
 * }}
 */
export function createProjectMemberApi(request) {
	const hide = { isHideTips: true };
	return {
		page: (body) => request(PROJECT_MEMBER_PAGE_API, { method: 'POST', ...hide, body }),
		candidate: (params) =>
			request(`${PROJECT_MEMBER_CANDIDATE_API}${buildProjectMemberCandidateQuery(params)}`, {
				method: 'GET',
				...hide,
			}),
		addMembers: (body) => request(PROJECT_MEMBER_API, { method: 'POST', ...hide, body }),
		updateRole: (id, body) =>
			request(getProjectMemberRoleApi(id), { method: 'PUT', ...hide, body }),
		delete: (id, body) =>
			request(getProjectMemberByIdApi(id), { method: 'DELETE', ...hide, body }),
	};
}
