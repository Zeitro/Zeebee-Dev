/**
 * 服务端取词：不依赖 next-intl 全局 plugin，供 RSC 在 ProjectIntlProvider 外使用。
 */
import { createTranslator } from 'next-intl';
import { PROJECT_LOCALE_COOKIE, resolveProjectLocale } from './config';
import { getProjectMessages } from './messages';

/**
 * @param {string | undefined | null} locale
 * @param {string} namespace
 */
export function getProjectTranslator(locale, namespace) {
	const resolved = resolveProjectLocale(locale);
	return createTranslator({
		locale: resolved,
		messages: getProjectMessages(resolved),
		namespace,
	});
}

/**
 * @param {import('next/headers').ReadonlyRequestCookies} cookieStore
 * @returns {'en' | 'zh'}
 */
export function readProjectLocaleFromCookies(cookieStore) {
	return resolveProjectLocale(cookieStore.get(PROJECT_LOCALE_COOKIE)?.value);
}

export { PROJECT_LOCALE_COOKIE };
