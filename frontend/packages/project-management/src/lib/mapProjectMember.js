/**
 * 项目成员 API 行 → 看板 UI 结构；补 `user_id` / `name` 兼容旧组件。
 */

/**
 * @param {object|null|undefined} item
 * @returns {object|null}
 */
export function mapProjectMember(item) {
	if (!item || typeof item !== 'object') return null;

	const organizationUserId = item.organization_user_id ?? item.user_id;

	return {
		...item,
		organization_user_id: organizationUserId,
		user_id: organizationUserId,
		full_name: item.full_name,
		email: item.email,
		name: item.full_name ?? item.name,
		role: item.role,
	};
}
