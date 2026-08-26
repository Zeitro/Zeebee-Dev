/**
 * 项目管理接口。分页与三端 FilterTable 协议一致：
 * POST `/project_management/project/page`，body `{ page, limit, filter }`，
 * 响应 data 为 `{ list, total }`（另含 page / limit / total_page）。
 * 创建：POST `/project_management/project`。
 * 详情：GET `/project_management/project/:id`，data 含 name / cover_url / workflow.stages。
 * 更新：PUT 同路径，body `{ name, description, cover_url }`。
 * 删除：DELETE 同路径。
 * 看板 / 成员仍走旧路径 `/projects/:id`，待对应接口迁移后再切。
 */

export const PROJECT_API = '/project_management/project';
export const PROJECT_PAGE_API = `${PROJECT_API}/page`;

const PROJECT_LEGACY_API = '/projects';

/**
 * 项目详情 / 更新 / 删除：GET / PUT / DELETE `/project_management/project/:id`。
 * @param {number|string} id
 * @returns {string}
 */
export function getProjectDetailApi(id) {
	return `${PROJECT_API}/${id}`;
}

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getProjectApi(id) {
	return `${PROJECT_LEGACY_API}/${id}`;
}

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getProjectMembersApi(id) {
	return `${PROJECT_LEGACY_API}/${id}/members`;
}

/**
 * @param {number|string} projectId
 * @param {number|string} userId
 * @returns {string}
 */
export function getProjectMemberApi(projectId, userId) {
	return `${PROJECT_LEGACY_API}/${projectId}/members/${userId}`;
}

/**
 * @param {number|string} projectId
 * @param {number|string} userId
 * @returns {string}
 */
export function getProjectMemberRoleApi(projectId, userId) {
	return `${PROJECT_LEGACY_API}/${projectId}/members/${userId}/role`;
}
