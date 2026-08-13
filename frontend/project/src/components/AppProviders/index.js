/**
 * 客户端全局 Provider：同步 antd / dayjs 与当前 next-intl locale。
 * 不改业务数据；仅 UI 本地化。
 */
'use client';

import { useEffect } from 'react';
import { ConfigProvider } from 'antd';
import { AntdRegistry } from '@ant-design/nextjs-registry';
import enUS from 'antd/locale/en_US';
import zhCN from 'antd/locale/zh_CN';
import dayjs from 'dayjs';
import 'dayjs/locale/en';
import 'dayjs/locale/zh-cn';

const antdLocales = {
	en: enUS,
	zh: zhCN,
};

/**
 * @param {{ locale: string, children: React.ReactNode }} props
 */
export default function AppProviders({ locale, children }) {
	useEffect(() => {
		dayjs.locale(locale === 'zh' ? 'zh-cn' : 'en');
	}, [locale]);

	return (
		<ConfigProvider locale={antdLocales[locale] || enUS}>
			<AntdRegistry>{children}</AntdRegistry>
		</ConfigProvider>
	);
}
