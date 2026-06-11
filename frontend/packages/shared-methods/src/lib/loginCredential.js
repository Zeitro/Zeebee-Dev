/**
 * 登录凭据：Remember me 仅存邮箱（localStorage）；密码由浏览器密码管理器保管。
 * 改密后通过 sessionStorage 一次性带到登录页（标签页关闭即失效）。
 */
import { localStorageUtil, sessionStorageUtil } from './local.js';

export const REMEMBERED_ACCOUNT_KEY = 'remembered_account';
const PENDING_LOGIN_CREDENTIAL_KEY = 'pending_login_credential';

/**
 * @returns {string|undefined}
 */
export function getRememberedAccountEmail() {
	return localStorageUtil.get(REMEMBERED_ACCOUNT_KEY)?.email;
}

/**
 * @param {string} email
 */
export function setRememberedAccountEmail(email) {
	if (email) {
		localStorageUtil.set(REMEMBERED_ACCOUNT_KEY, { email });
	}
}

export function clearRememberedAccount() {
	localStorageUtil.remove(REMEMBERED_ACCOUNT_KEY);
}

/**
 * 改密成功后写入，登录页消费后即删（仅当前标签页、一次性）。
 * @param {{ email: string, password: string }} credential
 */
export function setPendingLoginCredential({ email, password }) {
	if (!email || !password) return;
	sessionStorageUtil.set(PENDING_LOGIN_CREDENTIAL_KEY, { email, password });
}

/**
 * 读取并清除待登录凭据。
 * @returns {{ email: string, password: string }|null}
 */
export function consumePendingLoginCredential() {
	const credential = sessionStorageUtil.get(PENDING_LOGIN_CREDENTIAL_KEY);
	if (credential?.email && credential?.password) {
		sessionStorageUtil.remove(PENDING_LOGIN_CREDENTIAL_KEY);
		return credential;
	}
	return null;
}

/**
 * 从浏览器密码管理器静默读取已保存的登录密码。
 * @param {string} [email] 已记住的邮箱，用于校验凭据是否匹配
 * @returns {Promise<{ email: string, password: string }|null>}
 */
export async function getBrowserPasswordCredential(email) {
	if (typeof window === 'undefined' || !navigator.credentials?.get) return null;

	try {
		const credential = await navigator.credentials.get({
			password: true,
			mediation: 'silent',
		});
		if (
			credential?.type === 'password' &&
			credential.password &&
			credential.id &&
			(!email || credential.id === email)
		) {
			return { email: credential.id, password: credential.password };
		}
	} catch {
		// 无已保存凭据或用户未授权时静默降级
	}
	return null;
}

/**
 * 将登录凭据写入浏览器密码管理器；需在用户点击登录等手势回调中调用。
 * @param {string} email
 * @param {string} password
 * @returns {Promise<boolean>}
 */
export async function storeBrowserPasswordCredential(email, password) {
	if (typeof window === 'undefined' || !email || !password) return false;
	const PasswordCredential = window.PasswordCredential;
	if (!PasswordCredential || !navigator.credentials?.store) return false;

	try {
		const credential = new PasswordCredential({
			id: email,
			password,
			name: email,
		});
		await navigator.credentials.store(credential);
		return true;
	} catch {
		return false;
	}
}
