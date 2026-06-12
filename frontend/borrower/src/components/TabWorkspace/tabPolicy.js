import { matchRoute } from '@repo/shared-methods/utils';

/** 单 Tab 模式：侧边栏跳转时替换当前 Tab */
export const TAB_MODE_SINGLE = 'single';

/** 多 Tab 模式：loan 工作区内可并存多个 Tab（P1 起配合保活） */
export const TAB_MODE_MULTI = 'multi';

/**
 * loan 工作区路径：创建/编辑 loan 时可与 transaction 等页面多 Tab 并存。
 * P0 仅影响 Tab 开关策略，内容区仍走 Next children。
 */
export const LOAN_WORKSPACE_PATTERNS = [
	'/loans/createZeeBeeLoan',
	'/loans/zeeBeeDetails',
	'/myTransactions',
	'/myTransactions/:id',
	'/properties/:id',
];

/**
 * @param {string} pathname
 * @returns {{ mode: string, maxTabs: number }}
 */
export function getTabPolicy(pathname) {
	const isLoanWorkspace = LOAN_WORKSPACE_PATTERNS.some((pattern) =>
		matchRoute(pattern, pathname),
	);
	if (isLoanWorkspace) {
		return { mode: TAB_MODE_MULTI, maxTabs: 8 };
	}
	return { mode: TAB_MODE_SINGLE, maxTabs: 1 };
}
