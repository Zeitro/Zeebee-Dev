/**
 * 项目管理接口。分页与三端 FilterTable 协议一致：
 * POST body `{ page, limit, filter }`，响应 data 为 `{ list, total }`。
 */

export const PROJECT_PAGE_API = '/projects/page';
export const PROJECT_API = '/projects';
