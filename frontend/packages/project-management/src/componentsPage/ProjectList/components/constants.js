/**
 * 项目列表页内常量：创建卡占位不参与分页 total，仅插在当前页网格首位。
 */

export const CREATE_PROJECT_ITEM = Object.freeze({ __isCreateProject: true });

/**
 * @param {object|null|undefined} row
 * @returns {boolean}
 */
export function isCreateProjectItem(row) {
	return row?.__isCreateProject === true;
}
