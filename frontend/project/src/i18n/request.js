/**
 * next-intl 请求级配置：从 Cookie 读取 locale，并加载对应 messages。
 * 由 next.config 的 createNextIntlPlugin 指向本文件。
 */
import { cookies } from 'next/headers';
import { getRequestConfig } from 'next-intl/server';
import { defaultLocale, LOCALE_COOKIE, resolveLocale } from './config';

export default getRequestConfig(async () => {
	const store = await cookies();
	const locale = resolveLocale(store.get(LOCALE_COOKIE)?.value) || defaultLocale;
	const messages = (await import(`../../messages/${locale}.json`)).default;

	return {
		locale,
		messages,
	};
});
