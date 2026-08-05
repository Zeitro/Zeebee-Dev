/**
 * 列表表格分页：固定档位 + 首次访问按屏幕宽度初始化 + 用户偏好持久化。
 * 全站统一档位，不按单页布局计算精确行数。
 */

export const TABLE_PAGE_SIZE_OPTIONS = [10, 15, 20, 25, 30];
export const TABLE_PAGE_SIZE_MIN = 10;
export const TABLE_PAGE_SIZE_MAX = 30;
export const TABLE_PAGE_SIZE_STORAGE_KEY = 'zb_table_page_size';

/**
 * 读取用于分档的屏幕宽度（CSS 像素）。
 * 使用 screen.width 而非 innerWidth：不受 F12 面板、窗口缩放影响，更接近「显示器档位」。
 * @returns {number}
 */
export function getScreenWidthForPageSize() {
	if (typeof window === 'undefined') return 1920;
	return window.screen?.width || 1920;
}

/**
 * 按屏幕宽度映射默认每页条数档位（仅首次无偏好时调用）。
 *
 * | 屏幕宽度 (screen.width) | 默认条数 | 典型场景     |
 * |------------------------|---------|-------------|
 * | ≤ 1366                 | 10      | 小屏 / 1366 |
 * | 1367 ~ 1920            | 15      | 1080p       |
 * | 1921 ~ 2560            | 20      | 2K          |
 * | > 2560                 | 30      | 4K 及以上   |
 *
 * @param {number} [screenWidth]
 * @returns {number}
 */
export function getDefaultTablePageSize(screenWidth) {
	const width = screenWidth ?? getScreenWidthForPageSize();

	if (width <= 1366) return 10;
	if (width <= 1920) return 15;
	if (width <= 2560) return 20;
	return 30;
}

/**
 * @param {string} [storageKey]
 * @returns {number|null}
 */
export function getStoredTablePageSize(storageKey = TABLE_PAGE_SIZE_STORAGE_KEY) {
	if (typeof window === 'undefined') return null;

	try {
		const value = Number(window.localStorage.getItem(storageKey));
		return TABLE_PAGE_SIZE_OPTIONS.includes(value) ? value : null;
	} catch {
		return null;
	}
}

/**
 * @param {number} size
 * @param {string} [storageKey]
 */
export function setStoredTablePageSize(size, storageKey = TABLE_PAGE_SIZE_STORAGE_KEY) {
	if (typeof window === 'undefined' || !TABLE_PAGE_SIZE_OPTIONS.includes(size)) return;

	try {
		window.localStorage.setItem(storageKey, String(size));
	} catch {
		// ignore quota / privacy mode
	}
}

/**
 * 初始化 pageSize：
 * 1. 已有 localStorage → 直接用（含首次自动写入的档位、或用户手动改的）
 * 2. 无记录 → 按 screen.width 算档位并写入 localStorage，后续刷新不再重算
 *
 * @param {string} [storageKey]
 * @returns {number}
 */
export function getInitialTablePageSize(storageKey = TABLE_PAGE_SIZE_STORAGE_KEY) {
	const stored = getStoredTablePageSize(storageKey);
	if (stored) return stored;

	const defaultSize = getDefaultTablePageSize();
	setStoredTablePageSize(defaultSize, storageKey);
	return defaultSize;
}
