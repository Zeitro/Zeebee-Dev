/**
 * 项目模板（模板库 / 我的模板共用一套接口，按 `filter.template_type` 区分）。
 * 分页：POST `/project_management/template/page`，
 * body `{ page, limit, filter: { keyword, template_type }, sorter: { column, sort } }`，
 * 响应 data 为 `{ list, total, page, limit, total_page }`；list 项含 `stages`（按 `sort_order` 排）。
 * 详情：GET `/project_management/template/:id`，data 与 list 项同结构（含 `stages`）。
 * 删除：DELETE 同路径，data 为空对象。
 */

export const TEMPLATE_API = '/project_management/template';
export const TEMPLATE_PAGE_API = `${TEMPLATE_API}/page`;

/**
 * @param {number|string} id
 * @returns {string}
 */
export function getTemplateApi(id) {
	return `${TEMPLATE_API}/${id}`;
}
