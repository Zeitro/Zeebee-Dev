'use client';

import { useContext } from 'react';
import { useParams, usePathname, useSearchParams } from 'next/navigation';
import { TabRouteContext } from './TabRouteProvider';

/**
 * 统一读取路由参数：Tab 保活 Panel 内走 TabRouteContext，否则回退 Next 导航 hooks。
 * @returns {{ params: Record<string, string>, pathname: string, searchParams: URLSearchParams|import('next/dist/client/components/navigation').ReadonlyURLSearchParams }}
 */
export function useTabRoute() {
	const ctx = useContext(TabRouteContext);
	const nextParams = useParams();
	const nextPathname = usePathname();
	const nextSearchParams = useSearchParams();

	if (ctx) {
		return ctx;
	}

	return {
		params: nextParams ?? {},
		pathname: nextPathname,
		searchParams: nextSearchParams,
	};
}
