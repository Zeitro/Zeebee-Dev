/**
 * 工作流接口失败 message：i18n 文案 + 后端 error.message（若有）。
 */

/**
 * @param {(key: string) => string} t
 * @param {string} messageKey i18n key
 * @param {unknown} [error]
 * @returns {string}
 */
export function formatApiError(t, messageKey, error) {
	const base = t(messageKey);
	const detail = error?.message;
	return detail ? `${base} (${detail})` : base;
}
