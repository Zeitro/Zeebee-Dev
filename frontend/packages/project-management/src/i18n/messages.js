/**
 * 按 locale 加载项目管理文案；供 Server Component 与 Provider 共用。
 */
import enMessages from '@repo/project-management/messages/en.json';
import zhMessages from '@repo/project-management/messages/zh.json';
import { defaultLocale, resolveProjectLocale } from './config';

/** @type {Record<'en' | 'zh', typeof enMessages>} */
const messagesByLocale = {
	en: enMessages,
	zh: zhMessages,
};

/**
 * @param {string | undefined | null} locale
 * @returns {typeof enMessages}
 */
export function getProjectMessages(locale) {
	const resolved = resolveProjectLocale(locale);
	return messagesByLocale[resolved] ?? messagesByLocale[defaultLocale];
}
