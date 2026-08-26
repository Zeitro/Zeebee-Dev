/**
 * 项目管理接口。分页与三端 FilterTable 协议一致：
 * POST `/project_management/project/page`，body `{ page, limit, filter }`，
 * 响应 data 为 `{ list, total }`（另含 page / limit / total_page）。
 */

export const PROJECT_PAGE_API = '/project_management/project/page';
export const PROJECT_API = '/projects';

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getProjectApi(id) {
	return `${PROJECT_API}/${id}`;
}

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getProjectMembersApi(id) {
	return `${PROJECT_API}/${id}/members`;
}

/**
 * @param {number|string} projectId
 * @param {number|string} userId
 * @returns {string}
 */
export function getProjectMemberApi(projectId, userId) {
	return `${PROJECT_API}/${projectId}/members/${userId}`;
}

/**
 * @param {number|string} projectId
 * @param {number|string} userId
 * @returns {string}
 */
export function getProjectMemberRoleApi(projectId, userId) {
	return `${PROJECT_API}/${projectId}/members/${userId}/role`;
}
