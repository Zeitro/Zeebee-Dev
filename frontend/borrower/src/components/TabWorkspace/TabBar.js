'use client';

import React, { useLayoutEffect } from 'react';
import { Tabs } from 'antd';
import { usePathname, useRouter } from 'next/navigation';
import { eventBus } from '@repo/shared-methods/utils';
import { useUserStore } from '@repo/shared-methods/store';
import { useTabStore } from '@/store/tabStore';
import { parseTabKey, shouldShowTabBar } from './tabUtils';
import { useSyncUrlToTabs } from './useSyncUrlToTabs';

/**
 * 顶栏可编辑标签页（P0：切换 Tab = router.push，不保活页面实例）。
 */
const TabBar = () => {
	const router = useRouter();
	const pathname = usePathname();
	const { accessibleRoutes } = useUserStore();
	const { tabs, activeKey, closeTab, updateTabTitle } = useTabStore();

	useSyncUrlToTabs();

	const isShow = shouldShowTabBar(pathname, accessibleRoutes);

	// 兼容详情页 eventBus.changeCrumbs 动态标题
	useLayoutEffect(() => {
		const handleChangeCrumbs = (config = {}) => {
			const indices = Object.keys(config);
			if (!indices.length) return;
			const lastIndex = Math.max(...indices.map(Number));
			const title = config[lastIndex];
			if (title && activeKey) {
				updateTabTitle(activeKey, title);
			}
		};
		eventBus.on('changeCrumbs', handleChangeCrumbs);
		return () => {
			eventBus.off('changeCrumbs', handleChangeCrumbs);
		};
	}, [activeKey, updateTabTitle]);

	if (!isShow || !tabs.length) {
		return null;
	}

	const handleChange = (key) => {
		if (key === activeKey) return;
		const { pathname: nextPathname, search } = parseTabKey(key);
		const href = search ? `${nextPathname}?${search}` : nextPathname;
		// P0 内容仍走 Next 路由；syncTabFromUrl 会识别已存在 Tab 并保留其它标签
		router.push(href);
	};

	const handleEdit = (targetKey, action) => {
		if (action !== 'remove') return;
		const nextKey = closeTab(targetKey);
		if (!nextKey) return;
		const { pathname: nextPathname, search } = parseTabKey(nextKey);
		const href = search ? `${nextPathname}?${search}` : nextPathname;
		router.push(href);
	};

	return (
		<Tabs
			type="editable-card"
			hideAdd
			size="small"
			activeKey={activeKey}
			onChange={handleChange}
			onEdit={handleEdit}
			className="tab-workspace-bar mb-0 min-w-0 flex-1 [&_.ant-tabs-nav]:mb-0"
			items={tabs.map((tab) => ({
				key: tab.key,
				label: (
					<span className="max-w-[200px] truncate inline-block align-bottom">
						{tab.title}
					</span>
				),
				closable: tab.closable !== false && tabs.length > 1,
			}))}
		/>
	);
};

export default TabBar;
