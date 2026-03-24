import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { createRequest, secondsToDate } from '@repo/shared-methods/utils';
import Cookies from 'js-cookie';
const { request } = createRequest();
// 定义 cookie storage 适配器
const cookieStorage = {
	getItem: (name) => Cookies.get(name),
	setItem: (name, value) => {
		const jsonValue = JSON.parse(value || '{}');
		const tokenExpires = jsonValue?.state?.access?.access_token_expire; // 返回的秒
		const expires = secondsToDate(tokenExpires);
		Cookies.set(name, value, {
			path: '/',
			expires,
			secure: process.env.NEXT_PUBLIC_ENV === 'production',
			sameSite: 'strict',
		});
	},
	removeItem: (name) => Cookies.remove(name),
};
export const usePermissionStore = create(
	persist(
		(set, get, api) => ({
			access: {}, //token
			permissionTags: [], //权限keys
			accessPaths: [], //有权限的页面路由
			clearPermission: () => {
				api.persist.clearStorage();
				set({ access: null, permissionTags: null, accessPaths: null });
			},
			login: async (values, config = {}) => {
				const access = await request('/user/login', {
					method: 'POST',
					body: values,
					...config,
				});
				set({ access });
				return access;
			},
			updateAccess: (access) => {
				set({ access });
			},
			setPermissionTags: (permissionTags) => {
				set({ permissionTags });
			},
			setAccessPaths: (accessPaths) => {
				set({ accessPaths });
			},
		}),
		{
			// 存储唯一标识（必填）,区分系统类型主要因为本地调试cookie会互相覆盖
			name: `permissionStore_${process.env.NEXT_PUBLIC_SYSTEM_TYPE}`,
			storage: createJSONStorage(() => cookieStorage),
			partialize: (state) => ({
				access: state.access,
				permissionTags: state.permissionTags,
				accessPaths: state.accessPaths,
			}),
		},
	),
);
