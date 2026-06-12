'use client';
import { useCallback, useEffect, useState } from 'react';
import { Form } from 'antd';
import { localStorageUtil } from '../lib/local.js';
import {
	REMEMBERED_ACCOUNT_KEY,
	clearRememberedAccount,
	consumePendingLoginCredential,
	getBrowserPasswordCredential,
	setRememberedAccountEmail,
	storeBrowserPasswordCredential,
} from '../lib/loginCredential.js';

/**
 * 重置/修改密码后：写入 Remember me 邮箱，并同步浏览器密码管理器。
 * 应在用户点击提交后、长耗时 API 之前调用以保留手势上下文。
 * @param {string} email
 * @param {string} password
 */
export async function syncCredentialAfterPasswordReset(email, password) {
	if (!email || !password) return;
	setRememberedAccountEmail(email);
	await storeBrowserPasswordCredential(email, password);
}

/**
 * Remember me：仅持久化邮箱；密码从浏览器密码管理器读取，不写入 localStorage。
 * @param {import('antd').FormInstance} form
 * @returns {{ persistRememberedEmail: Function, shouldAutofill: boolean }}
 */
export function useRememberLoginEmail(form) {
	const remember = Form.useWatch('remember', form);
	const [hasRememberedEmail, setHasRememberedEmail] = useState(false);
	const shouldAutofill = hasRememberedEmail || remember === true;

	useEffect(() => {
		let cancelled = false;

		const hydrate = async () => {
			const pending = consumePendingLoginCredential();
			if (pending?.email && pending?.password) {
				setRememberedAccountEmail(pending.email);
				if (!cancelled) {
					form.setFieldsValue({
						email: pending.email,
						password: pending.password,
						remember: true,
					});
					setHasRememberedEmail(true);
				}
				return;
			}

			const account = localStorageUtil.get(REMEMBERED_ACCOUNT_KEY);
			if (account?.email) {
				const fields = { email: account.email, remember: true };
				const browserCred = await getBrowserPasswordCredential(account.email);
				if (!cancelled && browserCred?.password) {
					fields.password = browserCred.password;
				}
				if (!cancelled) {
					form.setFieldsValue(fields);
					setHasRememberedEmail(true);
				}
			} else if (!cancelled) {
				form.setFieldsValue({
					email: undefined,
					password: undefined,
					remember: false,
				});
				setHasRememberedEmail(false);
			}
		};

		hydrate();
		return () => {
			cancelled = true;
		};
	}, [form]);

	/** 仅持久化邮箱 */
	const persistRememberedEmail = useCallback((email, remember) => {
		if (remember && email) {
			setRememberedAccountEmail(email);
		} else {
			clearRememberedAccount();
		}
	}, []);

	return { persistRememberedEmail, shouldAutofill };
}
