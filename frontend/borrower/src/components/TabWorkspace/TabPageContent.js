'use client';

import React, { createElement, lazy, Suspense } from 'react';
import { Spin } from 'antd';
import { tabRegistry, resolveRegistryMatch } from './tabRegistry';

/** 模块加载时按 pattern 创建 lazy，避免 render 内调用 lazy() */
const lazyPageByPattern = new Map(
	tabRegistry.map(({ pattern, loader }) => [pattern, lazy(loader)]),
);

const PageFallback = () => (
	<div className="flex-1 flex items-center justify-center">
		<Spin />
	</div>
);

/**
 * 按 pathname 懒加载已注册的页面组件。
 * @param {{ pathname: string }} props
 */
const TabPageContent = ({ pathname }) => {
	const match = resolveRegistryMatch(pathname);
	if (!match) return null;

	const LazyPage = lazyPageByPattern.get(match.pattern);
	if (!LazyPage) return null;

	return <Suspense fallback={<PageFallback />}>{createElement(LazyPage)}</Suspense>;
};

export default TabPageContent;
