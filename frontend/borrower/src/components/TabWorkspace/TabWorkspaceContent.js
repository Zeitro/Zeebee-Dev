'use client';

import React, { useMemo } from 'react';
import { usePathname, useSearchParams } from 'next/navigation';
import { useTabStore } from '@/store/tabStore';
import { buildTabKey } from './tabUtils';
import { isRegisteredTabPath, resolveRegistryMatch } from './tabRegistry';
import { shouldKeepAliveTab } from './tabPolicy';
import TabRouteProvider from './TabRouteProvider';
import TabPanel from './TabPanel';
import TabPageContent from './TabPageContent';

/**
 * Tab 内容区：已注册路由走保活 Panel，未注册路由回退 Next {children}。
 * @param {{ children: React.ReactNode }} props
 */
const TabWorkspaceContent = ({ children }) => {
	const pathname = usePathname();
	const searchParams = useSearchParams();
	const { tabs, activeKey } = useTabStore();

	const activeTabKey = buildTabKey(pathname, searchParams);

	const cachedTabs = useMemo(
		() => tabs.filter((tab) => isRegisteredTabPath(tab.pathname)),
		[tabs],
	);

	const activeIsCached = useMemo(() => {
		if (!isRegisteredTabPath(pathname)) return false;
		return tabs.some((tab) => tab.key === activeTabKey);
	}, [pathname, activeTabKey, tabs]);

	return (
		<div className="flex-1 overflow-hidden flex flex-col relative">
			{cachedTabs.map((tab) => {
				const registryMatch = resolveRegistryMatch(tab.pathname);
				const tabWithParams = {
					...tab,
					params: tab.params ?? registryMatch?.params ?? {},
				};

				return (
					<TabPanel
						key={tab.key}
						isActive={tab.key === activeKey}
						keepAlive={shouldKeepAliveTab(tab.pathname)}
					>
						<TabRouteProvider tab={tabWithParams}>
							<TabPageContent pathname={tab.pathname} />
						</TabRouteProvider>
					</TabPanel>
				);
			})}
			{!activeIsCached ? (
				<div className="flex-1 overflow-hidden flex flex-col">{children}</div>
			) : null}
		</div>
	);
};

export default TabWorkspaceContent;
