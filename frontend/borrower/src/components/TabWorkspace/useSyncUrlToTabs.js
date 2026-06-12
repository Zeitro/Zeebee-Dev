'use client';

import { usePathname, useSearchParams } from 'next/navigation';
import { useEffect } from 'react';
import { useUserStore } from '@repo/shared-methods/store';
import { useTabStore } from '@/store/tabStore';
import { buildTabKey, resolveTabTitle } from './tabUtils';

/**
 * URL → tabStore 同步：路由变化时打开/激活对应 Tab。
 * P0 内容区仍只渲染 Next children，切换 Tab 会触发路由跳转。
 */
export function useSyncUrlToTabs() {
	const pathname = usePathname();
	const searchParams = useSearchParams();
	const { accessibleRoutes } = useUserStore();
	const syncTabFromUrl = useTabStore((state) => state.syncTabFromUrl);

	useEffect(() => {
		if (!accessibleRoutes?.length) return;

		const key = buildTabKey(pathname, searchParams);
		const title = resolveTabTitle(pathname, accessibleRoutes);

		syncTabFromUrl({
			key,
			pathname,
			search: searchParams.toString(),
			title,
		});
	}, [pathname, searchParams, accessibleRoutes, syncTabFromUrl]);
}
