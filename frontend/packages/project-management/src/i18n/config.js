/**
 * 项目管理域语言配置：仅 /projects 路由使用，与 borrower 全局英文隔离。
 * 新增语言：补 locales / localeLabels，并增加 messages/<locale>.json。
 */

/** @type {readonly ['en', 'zh']} */
export const locales = ['en', 'zh'];

/** @type {'en' | 'zh'} */
export const defaultLocale = 'en';

/** 与 setProjectLocaleAction / layout 读取共用的 Cookie 名 */
export const PROJECT_LOCALE_COOKIE = 'project_locale';

/** next-intl 日期格式时区，避免 Client/Server 环境差导致 markup mismatch */
export const PROJECT_TIMEZONE = 'UTC';

export const localeLabels = {
	en: 'English',
	zh: '中文',
};

/**
 * @param {string | undefined | null} value
 * @returns {'en' | 'zh'}
 */
export function resolveProjectLocale(value) {
	return locales.includes(value) ? value : defaultLocale;
}
