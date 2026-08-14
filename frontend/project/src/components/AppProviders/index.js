/**
 * 客户端全局 Provider：同步 antd / dayjs 与当前 next-intl locale；
 * 并向 platform-core FilterTable / FilterTableCard 注入应用侧 NewTable、NewPagination。
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
import { FilterTableProvider } from '@repo/platform-core';
import { NewTable } from '@/components/NewTable';
import { NewPagination } from '@/components/NewPagination';

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
			<AntdRegistry>
				<FilterTableProvider TableComponent={NewTable} PaginationComponent={NewPagination}>
					{children}
				</FilterTableProvider>
			</AntdRegistry>
		</ConfigProvider>
	);
}
