/**
 * 语言切换 Server Action：写入 Cookie 后由 Next 自动刷新 RSC。
 */
'use server';

import { cookies } from 'next/headers';
import { LOCALE_COOKIE, resolveLocale } from './config';

/**
 * @param {string} nextLocale
 */
export async function setLocaleAction(nextLocale) {
	const locale = resolveLocale(nextLocale);
	const store = await cookies();
	store.set(LOCALE_COOKIE, locale, {
		path: '/',
		maxAge: 60 * 60 * 24 * 365,
		sameSite: 'lax',
	});
}
