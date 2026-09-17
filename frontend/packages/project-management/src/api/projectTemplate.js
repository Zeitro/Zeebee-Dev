/**
 * 我的项目模板。另存仍走 `getProjectTemplateApi`（POST `/project_management/project/:id/template`）。
 * 列表：GET `/project_management/project_template/all`，data 为 `{ list }`（或直接数组）。
 * 预览：GET `/project_management/project_template/:id`。
 * 删除：DELETE 同路径。
 */

export const PROJECT_TEMPLATE_API = '/project_management/project_template';
export const PROJECT_TEMPLATE_ALL_API = `${PROJECT_TEMPLATE_API}/all`;

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getProjectTemplateDetailApi(id) {
	return `${PROJECT_TEMPLATE_API}/${id}`;
}
