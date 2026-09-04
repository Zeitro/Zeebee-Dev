'use client';
import { useUserStore, usePermissionStore } from '@repo/shared-methods/store';
import { matchRoute } from '@repo/shared-methods/utils';
import { useRouter } from 'next/navigation';
import { menuSetting } from '@config/routes';

const useData = () => {
	const router = useRouter();
	const { getUserInfo, getUserMenus, getOrganization, getCurrentPlan } = useUserStore();
	const { setPermissionTags, setAccessPaths, updateAccess } = usePermissionStore();

	// 拿到用户token后登录(获取用户信息/获取权限信息)
	const getUserPermissions = async (access, successCallback = () => {}) => {
		if (access) {
			updateAccess(access);
			// 	setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		}
		await getUserInfo();
		const result = await getUserMenus(menuSetting);
		await getOrganization();
		await getCurrentPlan();
		const accessPaths = result?.accessPaths ?? [];
		setPermissionTags(result?.permissionTags);
		setAccessPaths(accessPaths);
		successCallback?.(false);
		const hasDashboard = accessPaths.some((path) => matchRoute(path, '/dashboard'));
		router.push(hasDashboard ? '/dashboard' : '/');
	};
	return { getUserPermissions };
};
export default useData;
