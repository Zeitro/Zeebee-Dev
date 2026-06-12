import { flattenNavRouter, getCrumbs, getCurrentRouterItem } from '@repo/shared-methods/utils';

/**
 * Tab 唯一 key：pathname + search，与浏览器地址栏一致。
 * @param {string} pathname
 * @param {URLSearchParams|string} [searchParams]
 * @returns {string}
 */
export function buildTabKey(pathname, searchParams) {
	const search =
		typeof searchParams === 'string' ? searchParams : searchParams?.toString?.() || '';
	return search ? `${pathname}?${search}` : pathname;
}

/**
 * 从 pathname 解析 query（用于 router.push）。
 * @param {string} tabKey
 * @returns {{ pathname: string, search: string }}
 */
export function parseTabKey(tabKey) {
	const qIndex = tabKey.indexOf('?');
	if (qIndex === -1) {
		return { pathname: tabKey, search: '' };
	}
	return {
		pathname: tabKey.slice(0, qIndex),
		search: tabKey.slice(qIndex + 1),
	};
}

/**
 * 根据菜单配置解析 Tab 标题（默认取面包屑最后一级）。
 * @param {string} pathname
 * @param {Array} accessibleRoutes
 * @returns {string}
 */
export function resolveTabTitle(pathname, accessibleRoutes) {
	if (!accessibleRoutes?.length) return pathname;
	const flatRoutes = flattenNavRouter(accessibleRoutes, 'children') || [];
	const crumbs = getCrumbs(pathname, flatRoutes);
	const lastTitle = crumbs[crumbs.length - 1]?.title;
	if (lastTitle) return lastTitle;
	const item = getCurrentRouterItem(pathname, accessibleRoutes, 'children');
	return item?.name || pathname;
}

/**
 * 是否展示 Tab 栏（与原先面包屑 hideCrumbs 规则一致）。
 * @param {string} pathname
 * @param {Array} accessibleRoutes
 * @returns {boolean}
 */
export function shouldShowTabBar(pathname, accessibleRoutes) {
	if (!accessibleRoutes?.length) return false;
	const item = getCurrentRouterItem(pathname, accessibleRoutes, 'children');
	return Boolean(item && !item?.hideCrumbs);
}
