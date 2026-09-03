/**
 * 应用语言配置（无 URL 前缀，靠 Cookie 持久化）。
 * 新增语言：补 locales / localeLabels，并增加 messages/<locale>.json。
 */

export const locales = ['en', 'zh'];

/** @type {'en' | 'zh'} */
export const defaultLocale = 'en';

/** 与 next-intl request / setLocaleAction 共用的 Cookie 名 */
export const LOCALE_COOKIE = 'locale';

export const localeLabels = {
	en: 'English',
	zh: '中文',
};

/**
 * @param {string | undefined | null} value
 * @returns {'en' | 'zh'}
 */
export function resolveLocale(value) {
	return locales.includes(value) ? value : defaultLocale;
}
