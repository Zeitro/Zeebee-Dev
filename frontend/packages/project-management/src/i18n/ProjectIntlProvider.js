/**
 * 项目管理域 Intl Provider：嵌套 ConfigProvider，仅影响 /projects 子树；
 * 不包裹 AntdRegistry，避免与 borrower 根 Provider 冲突。
 */
'use client';

import { useEffect } from 'react';
import { ConfigProvider } from 'antd';
import enUS from 'antd/locale/en_US';
import zhCN from 'antd/locale/zh_CN';
import dayjs from 'dayjs';
import 'dayjs/locale/en';
import 'dayjs/locale/zh-cn';
import { NextIntlClientProvider } from 'next-intl';
import { PROJECT_TIMEZONE } from './config';

/** @type {Record<'en' | 'zh', typeof enUS>} */
const antdLocales = {
	en: enUS,
	zh: zhCN,
};

const validateMessagesByLocale = {
	en: {
		required: 'Please enter ${label}!',
	},
	zh: {
		required: '请输入${label}！',
	},
};

/**
 * @param {{ locale: 'en' | 'zh', messages: Record<string, unknown>, children: React.ReactNode }} props
 */
export default function ProjectIntlProvider({ locale, messages, children }) {
	useEffect(() => {
		dayjs.locale(locale === 'zh' ? 'zh-cn' : 'en');
	}, [locale]);

	return (
		<NextIntlClientProvider locale={locale} messages={messages} timeZone={PROJECT_TIMEZONE}>
			<ConfigProvider
				locale={antdLocales[locale] || enUS}
				form={{
					validateMessages:
						validateMessagesByLocale[locale] || validateMessagesByLocale.en,
				}}
			>
				{children}
			</ConfigProvider>
		</NextIntlClientProvider>
	);
}
