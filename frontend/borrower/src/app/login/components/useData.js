'use client';
import { useUserStore, usePermissionStore } from '@repo/shared-methods/store';
import { useRouter } from 'next/navigation';
const useData = () => {
	const router = useRouter();
	const { getUserInfo, getUserMenus } = useUserStore();
	const { setPermissionTags, setAccessPaths, updateAccess } = usePermissionStore();

	// 拿到用户token后登录(获取用户信息/获取权限信息)
	const getUserPermissions = async (access, successCallback = () => {}) => {
		await getUserInfo();
		const result = await getUserMenus();
		setPermissionTags(result?.permissionTags);
		setAccessPaths(result?.accessPaths);
		if (access) {
			updateAccess(access);
			// 	setCookie('token', access?.access_token, access?.access_token_expire * 1000);
		}
		successCallback?.(false);
		router.push('/');
	};
	return { getUserPermissions };
};
export default useData;
