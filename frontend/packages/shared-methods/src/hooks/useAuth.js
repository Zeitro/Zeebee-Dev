import { useCallback, useMemo } from 'react';
import { usePermissionStore } from '@repo/shared-methods/store';
export const useAuth = () => {
	const { permissionTags } = usePermissionStore();

	// 检查单个权限
	const hasPermission = useCallback(
		(permission) => {
			if (!permission) return true;
			// 开发环境豁免
			// if (process.env.NODE_ENV === 'development') return true;
			return permissionTags?.includes(permission) ?? false;
		},
		[permissionTags],
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

	const getActiveByPermission = useCallback(
		({ permissionList, render, ...options }) => {
			const active = {
				title: 'Action',
				dataIndex: 'operation',
				align: 'left',
				fixed: 'right',
				minWidth: 100,
				render,
				hidden: permissionList?.length ? !hasAnyPermission(permissionList) : false,
				...options,
			};
			return active;
		},
		[hasAnyPermission],
	);

	const auth = useMemo(
		() => ({
			permissionTags,
			hasPermission,
			hasAnyPermission,
			hasAllPermissions,
			getActiveByPermission,
		}),
		[permissionTags, hasPermission, hasAnyPermission, hasAllPermissions, getActiveByPermission],
	);

	return auth;
};

export default useAuth;
