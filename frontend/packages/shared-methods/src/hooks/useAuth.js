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

	// 通过按钮权限，获取table操作栏
	const getActiveByPermission = useCallback(
		({ permissionList, buttonOptions, render, ...options }) => {
			const minWidth = 100;
			const active = {
				title: 'Action',
				dataIndex: 'operation',
				align: 'left',
				fixed: 'right',
				width: 'max-content',
				minWidth,
				render,
				// 第一种传permissionList控制显示操作栏
				hidden: permissionList?.length ? !hasAnyPermission(permissionList) : false,
				...options,
			};
			// 设置按钮总宽度
			if (buttonOptions?.length) {
				// 第二种传buttonOptions控制显示操作栏，并设置操作栏总宽度
				const permissionArr = buttonOptions.map((item) => item?.permission);
				if (permissionArr?.length && !permissionList?.length) {
					active.hidden = !hasAnyPermission(permissionArr);
				}
				const widths = buttonOptions.map((item) => {
					if (item?.permission && hasPermission(item?.permission)) {
						return item.width;
					}
					return 0;
				});
				const totalWidth = widths.reduce((a, b) => a + b, 0);
				active.width = Math.max(minWidth, totalWidth);
			}
			return active;
		},
		[hasAnyPermission, hasPermission],
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
