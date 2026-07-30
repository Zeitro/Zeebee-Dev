'use client';
import React, { useEffect, useMemo, useState } from 'react';
import { ConfigProvider } from 'antd';
import { AntdRegistry } from '@ant-design/nextjs-registry';
import { CommonSidebar, NewHeader } from '@repo/shared-methods';
import HeaderCustomAction from '../HeaderCustomAction';
import { matchRoute, flattenNavRouter } from '@repo/shared-methods/utils';
import { usePathname } from 'next/navigation';
import { usePermissionStore, useUserStore } from '@repo/shared-methods/store';
import { useDeferredSyncEffect } from '@repo/shared-methods/hooks';
import useContent from './useContent';
import Image from 'next/image';
import { noTokenRouter } from '@config/routes';

const defaultToken = {
	colorPrimary: '#1677ff',
	colorTextBase: 'rgba(0, 0, 0, 0.88)',
	colorBgBase: '#fff',
	colorBorder: '#e6f7ff',
	colorError: '#FF4D4F',
	colorWarning: '#faad14',
	colorSuccess: '#00a86b',
	colorTextDisabled: '#d9d9d9',
};
function getCssVarValue(name, key) {
	if (typeof window !== 'undefined' && typeof document !== 'undefined') {
		return (
			getComputedStyle(document.documentElement).getPropertyValue(name).trim() ||
			defaultToken[key]
		);
	}
	return defaultToken[key];
}

/**
 * 自定义antd主题，区分有无（导航栏和面包屑）的页面
 */
export default function CustomizeAntdTheme({ children, fontFamily }) {
	const { accessibleRoutes, userInfo } = useUserStore();
	const { accessPaths, access } = usePermissionStore();
	const [antdToken, setAntdToken] = useState({});
	const pathname = usePathname();
	const { personalList, handleClickPersonalItem } = useContent();
	const validateMessages = {
		required: 'Please enter ${label}!',
	};
	// 没有导航栏/面包屑的纯页面（例如：login
	const isNoNav = useMemo(() => {
		return flattenNavRouter(noTokenRouter).some((item) => matchRoute(item.path, pathname));
	}, [pathname]);

	const isPermissionPage = useMemo(() => {
		return accessPaths?.some((path) => matchRoute(path, pathname));
	}, [pathname, accessPaths]);

	// 处理特殊情况（有导航栏的页面遇到token到期清空了，则重新登录）
	useEffect(() => {
		if ((isPermissionPage || !isNoNav) && (!access || !userInfo)) {
			window.location.replace('/login');
			console.log('---退出---');
		}
	}, [access, isPermissionPage, userInfo, isNoNav]);

	useDeferredSyncEffect({
		deferred: () => {
			setAntdToken({
				colorPrimary: getCssVarValue('--base-primary', 'colorPrimary'),
				// colorBgContainer: getCssVarValue('--background'),
				// colorBgContainerSelected: getCssVarValue('--blue-light-1'),
				colorTextBase: getCssVarValue('--foreground', 'colorTextBase'),
				colorBgBase: getCssVarValue('--background', 'colorBgBase'), // 页面背景色
				// colorTextSecondary: getCssVarValue('var(--black-light)'), // 次要文本颜色
				colorBorder: getCssVarValue('--base-border', 'colorBorder'), // 边框颜色
				colorError: getCssVarValue('--base-error', 'colorError'), // 错误状态颜色
				colorWarning: getCssVarValue('--base-warning', 'colorWarning'), // 警告状态颜色
				colorSuccess: getCssVarValue('--base-success', 'colorSuccess'), // 成功状态颜色
				colorTextDisabled: getCssVarValue('--base-disabled', 'colorTextDisabled'), // 禁用状态文本颜色
				// colorTextLabel: '',
				borderRadius: 8, // 圆角大小
				fontSize: 14, // 基础字体大小
				// fontSizeHeading: 16, // 标题字体大小
				// controlHeight: 40, // 控件高度（如按钮、输入框）
				fontFamily, // 字体
				// ... 其他 token ...
			});
		},
		deps: [],
		isUseMicrotask: true,
	});

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
						rowHoverBg: '#e6f4ff',
					},
					TreeSelect: {
						colorBorder: 'var(--base-disabled)',
					},
				},
			}}
			form={{ validateMessages }}
		>
			<AntdRegistry>
				{!isPermissionPage ? (
					children
				) : (
					<div className="h-full flex items-stretch customRoot">
						<CommonSidebar
							// routerList={navRouter}
							routerList={accessibleRoutes}
							personalList={personalList}
							handleClickPersonalItem={handleClickPersonalItem}
							renderLogo={() => (
								<Image
									src="/images/logo.svg"
									width={95}
									height={30}
									alt="ZeeBee Logo"
								/>
							)}
						/>
						<div className="flex-1 overflow-hidden flex flex-col">
							<NewHeader>
								<HeaderCustomAction />
							</NewHeader>
							{/* <div className="flex-1 overflow-hidden">{children}</div> */}
							{children}
						</div>
					</div>
				)}
			</AntdRegistry>
		</ConfigProvider>
	);
}
