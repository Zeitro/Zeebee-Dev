import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { createRequest } from '../lib/fetch.js';
import {
	createPermissionCookieStorage,
	getAccessTokenExpireSeconds,
	getAccessTokenExpireSecondsFromAccessCookie,
	getPermissionStorageKeys,
} from '../lib/permissionStorage.js';

const { request } = createRequest();
const storageKeys = getPermissionStorageKeys();

/** access Cookie：过期时间取自 persist 信封内 access_token_expire */
const accessCookieStorage = createPermissionCookieStorage(getAccessTokenExpireSeconds);

/** paths Cookie 内无 token，过期时间与 access Cookie 对齐 */
const pathsCookieStorage = createPermissionCookieStorage(
	(jsonValue) =>
		getAccessTokenExpireSeconds(jsonValue) ??
		getAccessTokenExpireSecondsFromAccessCookie(storageKeys.access),
);

/** token：Cookie，供 middleware / fetch 读取 */
export const usePermissionAccessStore = create(
	persist(
		(set) => ({
			access: {},
			updateAccess: (access) => set({ access }),
			clearAccess: () => {
				usePermissionAccessStore.persist.clearStorage();
				set({ access: null });
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
		}),
		{
			name: storageKeys.access,
			storage: createJSONStorage(() => accessCookieStorage),
			partialize: (state) => ({ access: state.access }),
			version: 1,
		},
	),
);

/** 可访问路由：Cookie，供 middleware 鉴权 */
export const usePermissionPathsStore = create(
	persist(
		(set) => ({
			accessPaths: [],
			setAccessPaths: (accessPaths) => set({ accessPaths: accessPaths ?? [] }),
			clearAccessPaths: () => {
				usePermissionPathsStore.persist.clearStorage();
				set({ accessPaths: [] });
			},
		}),
		{
			name: storageKeys.paths,
			storage: createJSONStorage(() => pathsCookieStorage),
			partialize: (state) => ({ accessPaths: state.accessPaths }),
			version: 1,
		},
	),
);

/** 按钮级权限 tag：localStorage，仅客户端 useAuth 使用 */
export const usePermissionTagsStore = create(
	persist(
		(set) => ({
			permissionTags: [],
			setPermissionTags: (permissionTags) => set({ permissionTags: permissionTags ?? [] }),
			clearPermissionTags: () => {
				usePermissionTagsStore.persist.clearStorage();
				set({ permissionTags: [] });
			},
		}),
		{
			name: storageKeys.tags,
			storage: createJSONStorage(() => localStorage),
			partialize: (state) => ({ permissionTags: state.permissionTags }),
			version: 1,
		},
	),
);

/** 登出 / 进入登录页：清空三个 store */
export function clearPermission() {
	usePermissionAccessStore.getState().clearAccess();
	usePermissionPathsStore.getState().clearAccessPaths();
	usePermissionTagsStore.getState().clearPermissionTags();
}

/**
 * 登录成功后一次性写入，避免分次 set 导致 access Cookie 中间态为 null。
 * @param {{ access?: object|null, permissionTags?: string[], accessPaths?: string[] }} payload
 */
export function setPermissions({ access, permissionTags, accessPaths }) {
	usePermissionAccessStore.getState().updateAccess(access ?? null);
	usePermissionPathsStore.getState().setAccessPaths(accessPaths ?? []);
	usePermissionTagsStore.getState().setPermissionTags(permissionTags ?? []);
}

/**
 * 兼容旧用法：聚合三个 store 的 state 与 actions。
 * 新代码可按需直接用 usePermissionAccessStore 等。
 */
export function usePermissionStore() {
	const access = usePermissionAccessStore((s) => s.access);
	const accessPaths = usePermissionPathsStore((s) => s.accessPaths);
	const permissionTags = usePermissionTagsStore((s) => s.permissionTags);

	return {
		access,
		accessPaths,
		permissionTags,
		login: usePermissionAccessStore.getState().login,
		updateAccess: usePermissionAccessStore.getState().updateAccess,
		setAccessPaths: usePermissionPathsStore.getState().setAccessPaths,
		setPermissionTags: usePermissionTagsStore.getState().setPermissionTags,
		setPermissions,
		clearPermission,
	};
}
