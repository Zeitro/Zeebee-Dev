/**
 * 任务类型颜色：ColorPicker 预设 / 表格 Tag；兼容命名色（blue）与 hex。
 */

export const TASK_TYPE_COLOR_HEX = {
	gray: '#8c8c8c',
	blue: '#1890ff',
	green: '#52c41a',
	orange: '#fa8c16',
	red: '#f5222d',
	purple: '#722ed1',
	cyan: '#13c2c2',
};

export const TASK_TYPE_DEFAULT_COLOR = TASK_TYPE_COLOR_HEX.blue;

export const TASK_TYPE_COLOR_PRESETS = Object.values(TASK_TYPE_COLOR_HEX);

const HEX_TO_NAME = Object.fromEntries(
	Object.entries(TASK_TYPE_COLOR_HEX).map(([name, hex]) => [hex.toLowerCase(), name]),
);

/**
 * 规范为 hex，供 ColorPicker / Tag；命名色走预设表。
 * @param {string|{toHexString?: () => string}|null|undefined} value
 * @returns {string}
 */
export function toTaskTypeColorHex(value) {
	if (!value) return '';
	if (typeof value !== 'string') {
		return value.toHexString?.() || '';
	}
	const trimmed = value.trim();
	if (trimmed.startsWith('#')) return trimmed;
	const name = trimmed === 'grey' ? 'gray' : trimmed;
	return TASK_TYPE_COLOR_HEX[name] || trimmed;
}

/**
 * Tag 文案：预设色用英文名（与图示一致），其余用 hex。
 * @param {string|{toHexString?: () => string}|null|undefined} value
 * @returns {string}
 */
export function getTaskTypeColorLabel(value) {
	const hex = toTaskTypeColorHex(value);
	if (!hex) return '';
	return HEX_TO_NAME[hex.toLowerCase()] || hex;
}

/**
 * 提交给创建 / 更新接口的 color：预设色用英文名（如 red），自定义色保留 hex。
 * @param {string|{toHexString?: () => string}|null|undefined} value
 * @returns {string}
 */
export function toTaskTypeColorPayload(value) {
	const hex = toTaskTypeColorHex(value);
	if (!hex) return '';
	return HEX_TO_NAME[hex.toLowerCase()] || hex;
}
