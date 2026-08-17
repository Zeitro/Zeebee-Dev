/**
 * 项目工作流：看板列来自 `GET /projects/:id/workflow` 的 states。
 */

import { PROJECT_API } from './project';

/**
 * @param {number|string} projectId
 * @returns {string}
 */
export function getProjectWorkflowApi(projectId) {
	return `${PROJECT_API}/${projectId}/workflow`;
}
