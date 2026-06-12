'use client';
import { useUserStore, usePermissionStore } from '@repo/shared-methods/store';
import { useRouter } from 'next/navigation';
import { menuSetting } from '@config/routes';
const useData = () => {
	const router = useRouter();
	const { getUserInfo, getUserMenus, getOrganization } = useUserStore();
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
		setPermissionTags(result?.permissionTags);
		setAccessPaths(result?.accessPaths);
		successCallback?.(false);
		router.push('/');
	};
	return { getUserPermissions };
};
export default useData;
