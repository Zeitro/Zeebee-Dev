// packages/shared-methods/src/hooks/useAuth.js
import { useCallback, useMemo } from 'react';
import { useUserStore } from '@/store/user';
export const useAuth = () => {
	const { permissions, loading } = useUserStore();

	// 检查单个权限
	const hasPermission = useCallback(
		(permission) => {
			if (!permission) return true;
			// 开发环境豁免
			if (process.env.NODE_ENV === 'development') return true;
			return permissions?.includes(permission) ?? false;
		},
		[permissions],
	);

	// 检查任一权限（满足一个即可）
	const hasAnyPermission = useCallback(
		(permissionList) => {
			if (!permissionList || !Array.isArray(permissionList)) return true;
			return permissionList.some((p) => hasPermission(p));
		},
		[hasPermission],
	);

	// 检查所有权限（必须全部满足）
	const hasAllPermissions = useCallback(
		(permissionList) => {
			if (!permissionList || !Array.isArray(permissionList)) return true;
			return permissionList.every((p) => hasPermission(p));
		},
		[hasPermission],
	);

	const auth = useMemo(
		() => ({
			permissions,
			loading,
			hasPermission,
			hasAnyPermission,
			hasAllPermissions,
		}),
		[permissions, loading, hasPermission, hasAnyPermission, hasAllPermissions],
	);

	return auth;
};

export default useAuth;
