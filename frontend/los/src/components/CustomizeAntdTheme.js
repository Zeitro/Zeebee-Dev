'use client';
import React, { useEffect, useMemo, useState } from 'react';
import { ConfigProvider } from 'antd';
import { AntdRegistry } from '@ant-design/nextjs-registry';
import SiderBar from '@/components/Sidebar';
import NewHeader from '@/components/NewHeader';
import { getCurrentRouterItem } from '@config/common';
import { noNavRouter } from '@config/routes';
import { usePathname } from 'next/navigation';
function getCssVarValue(name) {
	if (typeof window !== 'undefined' && typeof document !== 'undefined') {
		return getComputedStyle(document.documentElement).getPropertyValue(name).trim();
	}
	return '';
}

const defaultToken = {
	colorPrimary: '#1677ff',
	colorBgContainer: '#fff',
	colorBgContainerSelected: '#e6f7ff',
	colorTextBase: '#222',
	// ... 其他 token ...
};
/**
 * 自定义antd主题，区分有无（导航栏和面包屑）的页面
 */
export default function CustomizeAntdTheme({ children, fontFamily }) {
	const [antdToken, setAntdToken] = useState({});
	const pathname = usePathname();

    // 没有导航栏/面包屑的纯页面（例如：login
	const isNoNav = useMemo(() => {
		const routerItem = getCurrentRouterItem(pathname, noNavRouter);
		return !!routerItem;
	}, [pathname]);

	useEffect(
		() => {
			setAntdToken({
				colorPrimary: getCssVarValue('--base-primary'),
				colorBgContainer: getCssVarValue('--background'),
				colorBgContainerSelected: getCssVarValue('--blue-light-1'),
				colorTextBase: getCssVarValue('--foreground'),
				colorBgBase: getCssVarValue('var(--background)'), // 页面背景色
				colorTextSecondary: getCssVarValue('var(--black-light)'), // 次要文本颜色
				colorBorder: getCssVarValue('var(--base-border)'), // 边框颜色
				colorError: getCssVarValue('var(--base-error)'), // 错误状态颜色
				colorWarning: getCssVarValue('var(--base-warning)'), // 警告状态颜色
				colorSuccess: getCssVarValue('var(--base-success)'), // 成功状态颜色
				colorTextDisabled: getCssVarValue('var(--base-disabled)'), // 禁用状态文本颜色
				colorTextLabel: '',
				borderRadius: 8, // 圆角大小
				fontSize: 14, // 基础字体大小
				// fontSizeHeading: 16, // 标题字体大小
				// controlHeight: 40, // 控件高度（如按钮、输入框）
				fontFamily, // 字体
				// ... 其他 token ...
			});
		},
		[
			// 依赖主题切换变量
		],
	);

	return (
		<ConfigProvider
			theme={{
				token: antdToken,
				components: {
					Menu: {
						itemBg: 'var(--background)', // 普通项背景
						itemColor: '#474849', // 普通项文字
						itemSelectedBg: 'var(--background)', // 选中项背景
						itemSelectedColor: 'var(--base-primary)', // 选中项文字
						itemHoverBg: 'var(--blue-light-1)', // 悬停背景
						itemHoverColor: 'var(--base-primary)', // 悬停文字
					},
					Table: {
						border: '1px solid var(--base-border)',
						headerBg: 'var(--gray-light)',
						// headerColor: 'var(--foreground)',
						rowBg: 'var(--background)',
						rowHoverBg: 'var(--gray-light)',
					},
					TreeSelect:{
						colorBorder: 'var(--base-disabled)'
					}
				},
			}}
		>
			<AntdRegistry>
				{isNoNav ? (
					children
				) : (
					<div className="h-full flex items-stretch">
						<SiderBar />
						<div className="flex-1 overflow-hidden flex flex-col">
							<NewHeader />
							<div className='flex-1 overflow-hidden'>{children}</div>
						</div>
					</div>
				)}
			</AntdRegistry>
		</ConfigProvider>
	);
}
