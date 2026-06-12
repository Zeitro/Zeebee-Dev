'use client';

import React, { createContext, useMemo } from 'react';

/**
 * Tab 内页面的路由上下文，供 useTabRoute 在保活 Panel 中替代 useParams / useSearchParams。
 */
export const TabRouteContext = createContext(null);

/**
 * @param {{ tab: { pathname: string, search?: string, params?: Record<string, string> }, children: React.ReactNode }} props
 */
const TabRouteProvider = ({ tab, children }) => {
	const value = useMemo(() => {
		const searchParams = new URLSearchParams(tab.search || '');
		return {
			pathname: tab.pathname,
			params: tab.params || {},
			searchParams,
		};
	}, [tab.pathname, tab.search, tab.params]);

	return <TabRouteContext.Provider value={value}>{children}</TabRouteContext.Provider>;
};

export default TabRouteProvider;
