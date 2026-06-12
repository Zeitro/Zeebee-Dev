import { matchRoute } from '@repo/shared-methods/utils';

/**
 * Tab 保活页面注册表：pattern → 懒加载页面组件。
 * 更具体的路径写在前面，避免被宽泛 pattern 抢先匹配。
 */
export const tabRegistry = [
	{
		pattern: '/loans/createZeeBeeLoan',
		loader: () => import('@/app/loans/createZeeBeeLoan/page'),
	},
	{
		pattern: '/loans/createLoan',
		loader: () => import('@/app/loans/createLoan/page'),
	},
	{
		pattern: '/loans/zeeBeeDetails',
		loader: () => import('@/app/loans/zeeBeeDetails/page'),
	},
	{
		pattern: '/loans/:id',
		loader: () => import('@/app/loans/[id]/page'),
	},
	{
		pattern: '/loans',
		loader: () => import('@/app/loans/page'),
	},
	{
		pattern: '/myTransactions/:id',
		loader: () => import('@/app/myTransactions/[id]/page'),
	},
	{
		pattern: '/myTransactions',
		loader: () => import('@/app/myTransactions/page'),
	},
	{
		pattern: '/properties/:id',
		loader: () => import('@/app/properties/[id]/page'),
	},
	{
		pattern: '/properties',
		loader: () => import('@/app/properties/page'),
	},
	{
		pattern: '/management/borrower',
		loader: () => import('@/app/management/borrower/page'),
	},
	{
		pattern: '/management/borrowingEntity',
		loader: () => import('@/app/management/borrowingEntity/page'),
	},
	{
		pattern: '/management/contactLender',
		loader: () => import('@/app/management/contactLender/page'),
	},
];

/**
 * 从动态路由 pattern 解析 params。
 * @param {string} pattern 如 `/myTransactions/:id`
 * @param {string} pathname 如 `/myTransactions/5`
 * @returns {Record<string, string>}
 */
export function extractRouteParams(pattern, pathname) {
	const patternParts = pattern.split('/').filter(Boolean);
	const pathParts = pathname.split('/').filter(Boolean);
	const params = {};

	patternParts.forEach((segment, index) => {
		if (segment.startsWith(':')) {
			params[segment.slice(1)] = pathParts[index] ?? '';
		}
	});

	return params;
}

/**
 * 匹配 pathname 对应的注册项。
 * @param {string} pathname
 * @returns {{ pattern: string, loader: Function, params: Record<string, string> }|null}
 */
export function resolveRegistryMatch(pathname) {
	const entry = tabRegistry.find((item) => matchRoute(item.pattern, pathname));
	if (!entry) return null;
	return {
		...entry,
		params: extractRouteParams(entry.pattern, pathname),
	};
}

/**
 * pathname 是否已注册 Tab 保活。
 * @param {string} pathname
 * @returns {boolean}
 */
export function isRegisteredTabPath(pathname) {
	return Boolean(resolveRegistryMatch(pathname));
}
