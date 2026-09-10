/**
 * 项目管理语言切换 Server Action：写入 project_locale Cookie 后由 Next 刷新 RSC。
 */
'use server';

import { cookies } from 'next/headers';
import { PROJECT_LOCALE_COOKIE, resolveProjectLocale } from './config';

/**
 * @param {string} nextLocale
 */
export async function setProjectLocaleAction(nextLocale) {
	const locale = resolveProjectLocale(nextLocale);
	const store = await cookies();
	store.set(PROJECT_LOCALE_COOKIE, locale, {
		path: '/',
		maxAge: 60 * 60 * 24 * 365,
		sameSite: 'lax',
	});
}
