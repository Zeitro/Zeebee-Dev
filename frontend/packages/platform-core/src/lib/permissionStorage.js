import Cookies from 'js-cookie';
import { secondsToDate } from './date.js';

/**
 * 权限拆分存储 key：access / accessPaths → Cookie；permissionTags → localStorage。
 * @param {string} [systemType]
 */
export function getPermissionStorageKeys(systemType = process.env.NEXT_PUBLIC_SYSTEM_TYPE) {
	return {
		access: `permission_access_${systemType}`,
		paths: `permission_paths_${systemType}`,
		tags: `permission_tags_${systemType}`,
	};
}

/**
 * 从 zustand persist 信封解析 access_token_expire（秒）。
 * @param {object} [jsonValue]
 * @returns {number|undefined}
 */
export function getAccessTokenExpireSeconds(jsonValue) {
	const expire = jsonValue?.state?.access?.access_token_expire;
	if (expire == null) return undefined;
	const seconds = typeof expire === 'number' ? expire : Number(expire);
	return Number.isFinite(seconds) && seconds > 0 ? seconds : undefined;
}

/**
 * 从已写入的 access Cookie 读取 access_token_expire，供 paths 等与 token 过期对齐。
 * @param {string} [accessStorageKey]
 * @returns {number|undefined}
 */
export function getAccessTokenExpireSecondsFromAccessCookie(
	accessStorageKey = getPermissionStorageKeys().access,
) {
	try {
		return getAccessTokenExpireSeconds(JSON.parse(Cookies.get(accessStorageKey) || '{}'));
	} catch {
		return undefined;
	}
}

/**
 * zustand persist 用 Cookie 适配器：读写 + 按 token 过期时间设 expires。
 * @param {(jsonValue: object) => number|undefined} [getExpireSeconds] 默认从 persist 信封 state.access 取；paths 可改为读 access Cookie
 */
export function createPermissionCookieStorage(getExpireSeconds = getAccessTokenExpireSeconds) {
	return {
		getItem: (name) => Cookies.get(name) ?? null,
		setItem: (name, value) => {
			const jsonValue = JSON.parse(value || '{}');
			const expires = secondsToDate(getExpireSeconds(jsonValue));
			Cookies.set(name, value, {
				path: '/',
				expires,
				secure: process.env.NEXT_PUBLIC_ENV === 'production',
				sameSite: 'strict',
			});
		},
		removeItem: (name) => Cookies.remove(name, { path: '/' }),
	};
}

/**
 * middleware / fetch 读取拆分后的 Cookie（zustand persist 信封：state.xxx）。
 * @param {{ get?: (name: string) => { value?: string } | undefined }} cookieStore
 */
export function readPermissionFromCookieStore(cookieStore) {
	const keys = getPermissionStorageKeys();
	const getValue = (name) => cookieStore?.get?.(name)?.value ?? null;

	let access = null;
	const accessRaw = getValue(keys.access);
	if (accessRaw) {
		try {
			access = JSON.parse(accessRaw)?.state?.access ?? null;
		} catch {
			access = null;
		}
	}

	let accessPaths = [];
	const pathsRaw = getValue(keys.paths);
	if (pathsRaw) {
		try {
			accessPaths = JSON.parse(pathsRaw)?.state?.accessPaths ?? [];
		} catch {
			accessPaths = [];
		}
	}

	let permissionTags = [];
	if (typeof window !== 'undefined') {
		try {
			const tagsRaw = localStorage.getItem(keys.tags);
			if (tagsRaw) {
				permissionTags = JSON.parse(tagsRaw)?.state?.permissionTags ?? [];
			}
		} catch {
			permissionTags = [];
		}
	}

	return { access, accessPaths, permissionTags };
}

export function getAccessTokenFromStorage() {
	const { access } = readPermissionFromCookieStore({
		get: (name) => ({ value: Cookies.get(name) }),
	});
	return access?.access_token ?? null;
}
