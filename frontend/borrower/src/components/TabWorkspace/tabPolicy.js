import { matchRoute } from '@repo/shared-methods/utils';

/** 单 Tab 模式：新打开时替换全部 Tab（仅用于明确配置的路径） */
export const TAB_MODE_SINGLE = 'single';

/** 多 Tab 模式：可并存多个 Tab，切换已打开的 Tab 不会关闭其它 Tab */
export const TAB_MODE_MULTI = 'multi';

/** 强制单 Tab 的路径（如需要可在此配置） */
export const SINGLE_TAB_PATTERNS = [];

/**
 * loan 工作区：创建/编辑 loan 时与 transaction、property 详情等并存。
 */
export const LOAN_WORKSPACE_PATTERNS = [
	'/loans/createZeeBeeLoan',
	'/loans/zeeBeeDetails',
	'/loans/createLoan',
	'/loans/:id',
	'/myTransactions',
	'/myTransactions/:id',
	'/properties/:id',
];

/** 默认最多同时打开的 Tab 数 */
export const DEFAULT_MAX_TABS = 8;

/**
 * @param {string} pathname
 * @returns {{ mode: string, maxTabs: number }}
 */
export function getTabPolicy(pathname) {
	const isSingleTab = SINGLE_TAB_PATTERNS.some((pattern) => matchRoute(pattern, pathname));
	if (isSingleTab) {
		return { mode: TAB_MODE_SINGLE, maxTabs: 1 };
	}
	return { mode: TAB_MODE_MULTI, maxTabs: DEFAULT_MAX_TABS };
}

/**
 * 当前路径是否属于 loan 工作区（用于 P1 保活策略，P0 仅作标记）。
 * @param {string} pathname
 * @returns {boolean}
 */
export function isLoanWorkspacePath(pathname) {
	return LOAN_WORKSPACE_PATTERNS.some((pattern) => matchRoute(pattern, pathname));
}

/** Tab 保活路径前缀（含子路径） */
export const KEEP_ALIVE_PREFIXES = ['/loans', '/management'];

/**
 * 是否对 Tab 做 DOM 保活（关闭 Tab 前不卸载）。
 * @param {string} pathname
 * @returns {boolean}
 */
export function shouldKeepAliveTab(pathname) {
	if (isLoanWorkspacePath(pathname)) return true;
	return KEEP_ALIVE_PREFIXES.some(
		(prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`),
	);
}
