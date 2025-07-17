'use client';
import React, { useEffect, useState } from 'react';
import { ConfigProvider } from 'antd';
import { AntdRegistry } from '@ant-design/nextjs-registry';

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

export default function ClientAntdProvider({ children }) {
	const [antdToken, setAntdToken] = useState({});

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
				fontFamily: 'Arial, Helvetica, sans-serif', // 字体
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
				},
			}}
		>
			<AntdRegistry>{children}</AntdRegistry>
		</ConfigProvider>
	);
}
