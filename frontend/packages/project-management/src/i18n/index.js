/**
 * 项目管理域 i18n 入口：Scoped Provider + 文案加载 + Server Action。
 * borrower 在 app/projects/layout.js 包 ProjectIntlProvider 即可，勿在根 layout 启用。
 */
export {
	locales,
	defaultLocale,
	PROJECT_LOCALE_COOKIE,
	localeLabels,
	resolveProjectLocale,
} from './config';
export { getProjectMessages } from './messages';
export { setProjectLocaleAction } from './actions';
export { getProjectTranslator, readProjectLocaleFromCookies } from './server';
export { default as ProjectIntlProvider } from './ProjectIntlProvider';
export { default as LocaleSwitcher } from './LocaleSwitcher';
